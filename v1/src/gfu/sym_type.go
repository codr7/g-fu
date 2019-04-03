package gfu

import (
  "strings"
)

type SymType struct {
  BasicType
}

func (t *SymType) Init(id *Sym) *SymType {
  t.BasicType.Init(id)
  return t
}

func (t *SymType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('\'')
  out.WriteString(val.AsSym().name)
}

func (v Val) AsSym() *Sym {
  return v.imp.(*Sym)
}
