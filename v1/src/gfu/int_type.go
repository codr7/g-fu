package gfu

type IntType struct {
  BasicType
}

func (t *IntType) Bool(g *G, val Val) bool {
  return val.AsInt() > 0
}

func (t *IntType) Splat(g *G, pos Pos, val Val, out []Val) []Val {
  for i := 0; i < val.AsInt(); i++ {
    var iv Val
    iv.Init(pos, t, i)
    out = append(out, iv)
  }

  return out
}

func (v Val) AsInt() int {
  return v.imp.(int)
}
