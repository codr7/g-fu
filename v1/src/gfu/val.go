package gfu

import (
  //"log"
  "strings"
)

type Val interface {
  Bool(*G) bool
  Call(*G, Vec, *Env) (Val, E)
  Dump(*strings.Builder)
  Eq(*G, Val) bool
  Eval(*G, *Env) (Val, E)
  Is(*G, Val) bool
  Quote(*G, *Env) (Val, E)
  Splat(*G, Vec) Vec
  Type(*G) *Type
}

func (env *Env) AddVal(g *G, id string, val Val) {
  env.Put(g.Sym(id), val)
}
