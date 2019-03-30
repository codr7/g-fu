package gfu

type PrimImp func (*G, ListForm, *Env, Pos) (Val, Error)

type Prim struct {
  id *Sym
  nargs int
  imp PrimImp
}

func NewPrim(id *Sym, nargs int, imp PrimImp) *Prim {
  p := new(Prim)
  p.id = id
  p.nargs = nargs
  p.imp = imp
  return p
}

func (e *Env) AddPrim(g *G, id *Sym, nargs int, imp PrimImp) {
  var p Val
  p.Init(g.Prim, NewPrim(id, nargs, imp))
  e.Put(id, p)
}
