package gfu

type IntType struct {
  BasicType
}

func (t *IntType) Bool(g *G, val Val) bool {
  return val.AsInt() > 0
}

func (v Val) AsInt() int {
  return v.imp.(int)
}
