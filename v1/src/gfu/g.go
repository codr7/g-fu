package gfu

type Syms map[string]*Sym

type G struct {
  Debug bool
  Pos Pos
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  Fun, Int Type
}

func NewG() (*G, Error) {
  return new(G).Init()
}

func (g *G) Init() (*G, Error) {
  g.syms = make(Syms)
  g.Pos.Init("n/a", -1, -1)
  g.Fun = new(FunType).Init(g.Sym("Fun"))
  g.Int = new(IntType).Init(g.Sym("Int"))
  return g, nil
}
