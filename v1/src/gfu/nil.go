package gfu

import (
  //"log"
  "strings"
)

type Nil struct {
  BasicVal
}

func (n *Nil) Init(g *G) *Nil {
  n.BasicVal.Init(&g.NilType, n)
  return n
}

func (_ *Nil) Bool(g *G) bool {
  return false
}

func (_ *Nil) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E("Nil call")
}

func (_ *Nil) Dump(out *strings.Builder) {
  out.WriteRune('_')
}

func (_ *Nil) Len(g *G) (Int, E) {
  return 0, nil
}

func (n *Nil) Pop(g *G) (Val, Val, E) {
  return n, n, nil
}

func (_ *Nil) Push(g *G, its...Val) (Val, E) {
  return Vec(its), nil
}

func (_ *Nil) Splat(g *G, out Vec) Vec {
  return out
}

