package gfu

type PrimImp func (*G, Pos, VecForm, *Env) (Val, E)

type Prim struct {
  id *Sym
  imp PrimImp
}

func NewPrim(id *Sym, imp PrimImp) *Prim {
  p := new(Prim)
  p.id = id
  p.imp = imp
  return p
}

func (e *Env) AddPrim(g *G, id string, imp PrimImp) {
  var p Val
  ids := g.S(id)
  p.Init(g.Prim, NewPrim(ids, imp))
  e.Put(ids, p)
}
