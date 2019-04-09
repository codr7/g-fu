package gfu

import (
  //"log"
  "strings"
)

type SplatType struct {
  BasicType
}

func (t *SplatType) Bool(g *G, val Val) bool {
  return val.AsSplat().AsBool(g)
}

func (t *SplatType) Dump(val Val, out *strings.Builder) {
  val.AsSplat().Dump(out)
  out.WriteString("..")
}

func (t *SplatType) Eq(g *G, x Val, y Val) bool {
  return x.AsSplat().Is(g, y.AsSplat())
}

func (v Val) AsSplat() Val {
  return v.imp.(Val)
}
