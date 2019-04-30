package gfu

func (g *G) Bool(val bool) Val {
  if val {
    return &g.T
  }

  return &g.F
}
