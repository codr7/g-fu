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

func (t *SymType) Unquote(g *G, pos Pos, val Val) (Form, E) {
  return new(IdForm).Init(pos, val.AsSym()), nil
}

func (v Val) AsSym() *Sym {
  return v.imp.(*Sym)
}
