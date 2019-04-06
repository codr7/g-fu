package gfu

type IntType struct {
  BasicType
}

func Abs(x int) int {
  if x < 0 {
    return -x
  }

  return x
}

func (t *IntType) Bool(g *G, val Val) bool {
  return val.AsInt() > 0
}

func (t *IntType) Splat(g *G, val Val, out []Val) []Val {
  for i := 0; i < val.AsInt(); i++ {
    var iv Val
    iv.Init(t, i)
    out = append(out, iv)
  }

  return out
}

func (v Val) AsInt() int {
  return v.imp.(int)
}
