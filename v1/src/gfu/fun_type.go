package gfu

import (
  "strings"
)

type FunType struct {
  BasicType
}

func (t *FunType) Init(id *Sym) *FunType {
  t.BasicType.Init(id)
  return t
}

func (t *FunType) Dump(val interface{}, out *strings.Builder) {
  f := val.(*Fun)
  out.WriteString("(fun (")

  for i, a := range f.args {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.name)
  }

  out.WriteString(") ")

  for i, bf := range f.body {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(bf.String())   
  }
  
  out.WriteRune(')')
}

func (v *Val) Fun() *Fun {
  return v.imp.(*Fun)
}
