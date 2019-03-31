package gfu

type Sym struct {
  tag Tag
  name string
}

func NewSym(g *G, name string) *Sym {
  return new(Sym).Init(g, name)
}

func (s *Sym) Init(g *G, name string) *Sym {
  s.tag = g.NextSymTag()
  s.name = name
  g.syms[name] = s
  return s
}

func (s *Sym) String() string {
  return s.name
}

func (g *G) S(name string) *Sym {
  if s := g.syms[name]; s != nil { return s }
  return new(Sym).Init(g, name)
}

func (g *G) NextSymTag() Tag {
  g.sym_tag++
  return g.sym_tag
}

