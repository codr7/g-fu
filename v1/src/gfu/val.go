package gfu

import (
  //"log"
  "strings"
)

type Val interface {
  Bool(*G) bool
  Call(*G, *Task, *Env, Vec) (Val, E)
  Dump(*strings.Builder)
  Eq(*G, Val) bool
  Eval(*G, *Task, *Env) (Val, E)
  Is(*G, Val) bool
  Quote(*G, *Task, *Env) (Val, E)
  Splat(*G, Vec) Vec
  Type(*G) *Type
}

func (env *Env) AddVal(g *G, id string, val Val) {
  env.Let(g.Sym(id), val)
}
