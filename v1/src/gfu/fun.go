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
