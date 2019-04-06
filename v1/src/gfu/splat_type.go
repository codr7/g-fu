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
  return x.AsSplat().Eq(g, y.AsSplat())
}

func (t *SplatType) Unquote(g *G, pos Pos, val Val) (Form, E) {
  f, e := val.AsSplat().Unquote(g, pos)

  if e != nil {
    return nil, e
  }
    
  return new(SplatForm).Init(pos, f), nil
}

func (v Val) AsSplat() Val {
  return v.imp.(Val)
}
