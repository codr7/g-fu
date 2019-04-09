package gfu

import (
  //"log"
  "strings"
)

type SpliceType struct {
  BasicType
}

func (t *SpliceType) Bool(g *G, val Val) bool {
  return val.AsSplice().AsBool(g)
}

func (t *SpliceType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('%')
  val.AsSplice().Dump(out)
}

func (t *SpliceType) Eq(g *G, x Val, y Val) bool {
  return x.AsSplice().Is(g, y.AsSplice())
}

func (t *SpliceType) Eval(g *G, pos Pos, val Val, env *Env) (Val, E) {
  return g.NIL, g.E(pos, "Unquoted splice")
}

func (t *SpliceType) Quote(g *G, pos Pos, val Val, env *Env) (Val, E) {
  return val.AsSplice().Eval(g, pos, env)
}

func (v Val) AsSplice() Val {
  return v.imp.(Val)
}
