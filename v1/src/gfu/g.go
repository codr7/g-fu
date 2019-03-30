package gfu

type Syms map[string]*Sym

type G struct {
  Debug bool
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  Fun, Int, Nil, Prim Type
  NIL Val
}

func NewG() (*G, Error) {
  return new(G).Init()
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

func (g *G) AddPrim(id *Sym, nargs int, imp PrimImp) {
  g.RootEnv.Put(id, g.Prim, NewPrim(id, nargs, imp))
}

func (g *G) Init() (*G, Error) {
  g.syms = make(Syms)
  g.Fun = new(FunType).Init(g.Sym("Fun"))
  g.Int = new(IntType).Init(g.Sym("Int"))
  g.Nil = new(NilType).Init(g.Sym("Nil"))
  g.Prim = new(PrimType).Init(g.Sym("Prim"))
  g.NIL.Init(g.Nil, nil)
  g.RootEnv.Put(g.Sym("_"), g.Nil, g.NIL)
  g.AddPrim(g.Sym("+"), 2, int_add_imp)
  return g, nil
}
