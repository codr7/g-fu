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

func (t *NilType) AsBool(g *G, _ Val) bool {
  return false
}

func (t *NilType) Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error) {
  return g.NIL, g.NewError(pos, "Nil call")
}
  
func (t *NilType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('_')
}
