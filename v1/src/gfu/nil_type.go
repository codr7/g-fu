package gfu

import (
  //"log"
  "strings"
)

type NilType struct {
  BasicType
}

func (t *NilType) Init(id *Sym) *NilType {
  t.BasicType.Init(id)
  return t
}

func (t *NilType) AsBool(g *G, val Val) bool {
  return false
}

func (t *NilType) Call(g *G, pos Pos, val Val, args VecForm, env *Env) (Val, E) {
  return g.NIL, g.E(pos, "Nil call")
}
  
func (t *NilType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('_')
}

func (t *NilType) Splat(g *G, val Val, out []Val) []Val {
  return out
}

