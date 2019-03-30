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

func (g *G) AddPrim(id *Sym, nargs int, imp PrimImp) {
  g.RootEnv.Put(id, g.Prim, NewPrim(id, nargs, imp))
}

func int_add_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  avs, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  var out Val
  out.Init(g.Int, avs[0].AsInt() + avs[1].AsInt())
  return out, nil
}

func (g *G) InitPrims() {
  g.AddPrim(g.Sym("+"), 2, int_add_imp)
}
