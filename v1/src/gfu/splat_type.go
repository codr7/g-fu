package gfu

import (
  //"log"
  "strings"
)

type SplatType struct {
  BasicType
}

func (t *SplatType) Init(id *Sym) *SplatType {
  t.BasicType.Init(id)
  return t
}

func (t *SplatType) AsBool(g *G, val Val) bool {
  return val.imp.(Val).AsBool(g)
}

func (t *SplatType) Dump(val Val, out *strings.Builder) {
  v := val.imp.(Val);
  v.Dump(out)
  out.WriteString("..")
}
