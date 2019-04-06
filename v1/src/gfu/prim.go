package gfu

type PrimImp func (*G, Pos, []Form, *Env) (Val, E)

type Prim struct {
  id *Sym
  arg_list ArgList
  imp PrimImp
}

func NewPrim(g *G, id *Sym, imp PrimImp, args []*Sym) *Prim {
  p := new(Prim)
  p.id = id
  p.arg_list.Init(g, args)
  p.imp = imp
  return p
}

func (e *Env) AddPrim(g *G, id string, imp PrimImp, args...string) {
  var p Val
  ids := g.Sym(id)
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.Sym(a)
  }

  p.Init(g.PrimType, NewPrim(g, ids, imp, as))
  e.Put(ids, p)
}
