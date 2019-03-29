package gfu

type IntType struct {
  BasicType
}

type Int int64

func (t *IntType) Init(id *Sym) *IntType {
  t.BasicType.Init(id)
  return t
}

func (v *Val) Int() Int {
  return v.imp.(Int)
}
