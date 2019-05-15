package gfu

type Num interface {
  Type
  Abs(*G, Val) (Val, E)
  Add(*G, Val, Val) (Val, E)
  Dec(*G, Val) (Val, E)
  Div(*G, Val, Val) (Val, E)
  Int(*G, Val) (Val, E)
  Mul(*G, Val, Val) (Val, E)
  Neg(*G, Val) (Val, E)
  Sub(*G, Val, Val) (Val, E)
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

func (g *G) Dec(x Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Dec(g, x)
}

func (g *G) Div(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Div(g, x, y)
}

func (g *G) Int(x Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Int(g, x)
}

func (g *G) Mul(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Mul(g, x, y)
}

func (g *G) Neg(x Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Neg(g, x)
}

func (g *G) Sub(x, y Val) (Val, E) {
  t := x.Type(g)
  nt, ok := t.(Num)

  if !ok {
    return nil, g.E("Expected Num: %v", t)
  }

  return nt.Sub(g, x, y)
}
