package gfu

type Syms map[string]*Sym

type G struct {
  Debug bool
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  Bool, Fun, Int, Meta, Nil, Prim, Splat, Vec Type
  NIL, T, F Val
}

func NewG() (*G, Error) {
  return new(G).Init()
}

func (g *G) Init() (*G, Error) {
  g.syms = make(Syms)
  return g, nil
}
