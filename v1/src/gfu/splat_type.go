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

func (t *SplatType) Eval(g *G, val Val, env *Env) (Val, E) {
  var e E
  val.imp, e = val.AsSplat().Eval(g, env)
  return val, e
}

func (t *SplatType) Quote(g *G, val Val, env *Env) (Val, E) {
  var e E

  if val, e = val.AsSplat().Quote(g, env); e != nil {
    return g.NIL, e
  }

  val.Init(g.SplatType, val)
  return val, nil
}

func (t *SplatType) Splat(g *G, val Val, out Vec) Vec {
  v := val.AsSplat()

  if v.val_type != g.VecType {
    return append(out, val)
  }

  return v.Splat(g, out)
}

func (v Val) AsSplat() Val {
  return v.imp.(Val)
}
