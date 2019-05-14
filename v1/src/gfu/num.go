package gfu

type Num interface {
  Type
  Abs(*G, Val) (Val, E)
  Add(*G, Val, Val) (Val, E)
  Sub(*G, Val, Val) (Val, E)
  Mul(*G, Val, Val) (Val, E)
  Div(*G, Val, Val) (Val, E)
}

type NumType struct {
  BasicType
}

func (g *G) Abs(x Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Abs(g, x)
}

func (g *G) Add(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Add(g, x, y)
}

func (g *G) Sub(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Sub(g, x, y)
}

func (g *G) Mul(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Mul(g, x, y)
}

func (g *G) Div(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Div(g, x, y)
}
