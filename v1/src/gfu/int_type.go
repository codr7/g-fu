package gfu

type IntType struct {
  BasicType
}

type Int int64

func (i Int) Abs() Int {
  if i < 0 {
    return -i
  }

  return i
}

func (t *IntType) Init(id *Sym) *IntType {
  t.BasicType.Init(id)
  return t
}

func (t *IntType) AsBool(g *G, val Val) bool {
  return val.AsInt() > 0
}

func (t *IntType) Splat(g *G, val Val, out []Val) []Val {
  for i := Int(0); i < val.AsInt(); i++ {
    var iv Val
    iv.Init(t, i)
    out = append(out, iv)
  }

  return out
}

func (v Val) AsInt() Int {
  return v.imp.(Int)
}
