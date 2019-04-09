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

func (t *SplatType) Quote(g *G, pos Pos, val Val, env *Env) (Val, E) {
  var e E

  if val, e = val.AsSplat().Quote(g, pos, env); e != nil {
    return g.NIL, e
  }

  val.Init(pos, g.SplatType, val)
  return val, nil
}

func (v Val) AsSplat() Val {
  return v.imp.(Val)
}
