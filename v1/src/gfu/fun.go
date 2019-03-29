package gfu

type Fun struct {
  args []*Sym
  body []Form
  env *Env
}

func NewFun(args []*Sym, body []Form, env *Env) *Fun {
  return new(Fun).Init(args, body, env)
}

func (f *Fun) Init(args []*Sym, body []Form, env *Env) *Fun {
  f.args = args
  f.body = body
  f.env = env
  return f
}

func (f *Fun) Call(g *G, vals []Val, env *Env) (*Val, Error) {
  if len(vals) != len(f.args) {
    return nil, g.NewError(&g.Pos, "Arg mismatch: %v", vals)
  }
  
  var out *Val
  var be Env
  f.env.Clone(&be)
  be.Merge(f.args, vals)
  
  for _, bf := range f.body {
    var e Error
    
    if out, e = bf.Eval(g, &be); e != nil {
      return nil, g.NewError(&g.Pos, "Call failed: %v", e)
    }
  }
  
  return out, nil
}
