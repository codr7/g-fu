package gfu

type PrimImp func (*G, Pos, VecForm, *Env) (Val, E)

type Prim struct {
  id *Sym
  min_args, max_args int
  imp PrimImp
}

func NewPrim(id *Sym, min_args, max_args int, imp PrimImp) *Prim {
  p := new(Prim)
  p.id = id
  p.min_args, p.max_args = min_args, max_args
  p.imp = imp
  return p
}

func (p *Prim) CheckArgs(g *G, pos Pos, args []Val) E {
  nargs := len(args)
  
  if (p.min_args != -1 && nargs < p.min_args) ||
    (p.max_args != -1 && nargs > p.max_args) {
    return g.E(pos, "Arg mismatch: %v", p.id)
  }

  return nil
}

func (e *Env) AddPrim(g *G, id *Sym, min_args, max_args int, imp PrimImp) {
  var p Val
  p.Init(g.Prim, NewPrim(id, min_args, max_args, imp))
  e.Put(id, p)
}
