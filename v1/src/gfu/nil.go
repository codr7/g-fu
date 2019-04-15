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

func (n *Nil) Eq(g *G, rhs Val) bool {
  return n == rhs
}

func (n *Nil) Eval(g *G, task *Task, env *Env) (Val, E) {
  return n, nil
}

func (n *Nil) Is(g *G, rhs Val) bool {
  return n == rhs
}

func (n *Nil) Quote(g *G, task *Task, env *Env) (Val, E) {
  return n, nil
}

func (_ *Nil) Splat(g *G, out Vec) Vec {
  return out
}

func (n *Nil) String() string {
  return "_"
}

