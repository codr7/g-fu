package gfu

import (
  //"log"
  "strings"
)

type Val interface {
  Type(*G) Type
}

func (g *G) Bool(val Val) (bool, E) {
  return val.Type(g).Bool(g, val)
}

func (g *G) Call(task *Task, env *Env, val Val, args Vec, args_env *Env) (Val, E) {
  return val.Type(g).Call(g, task, env, val, args, args_env)
}

func (g *G) Clone(val Val) (Val, E) {
  return val.Type(g).Clone(g, val)
}

func (g *G) Drop(val Val, n Int) (out Val, e E) {
  return val.Type(g).Drop(g, val, n)
}

func (g *G) Dup(val Val) (Val, E) {
  return val.Type(g).Dup(g, val)
}

func (g *G) Dump(val Val, out *strings.Builder) E {
  return val.Type(g).Dump(g, val, out)
}

func (g *G) Eq(lhs, rhs Val) (bool, E) {
  return lhs.Type(g).Eq(g, lhs, rhs)
}

func (g *G) Eval(task *Task, env *Env, val Val) (Val, E) {
  return val.Type(g).Eval(g, task, env, val)
}

func (g *G) Expand(task *Task, env *Env, val Val, depth Int) (Val, E) {
  return val.Type(g).Expand(g, task, env, val, depth)
}

func (g *G) Extenv(src, dst *Env, val Val, clone bool) E {  
  return val.Type(g).Extenv(g, src, dst, val, clone)
}

func (g *G) Is(lhs, rhs Val) bool {
  return lhs.Type(g).Is(g, lhs, rhs)
}

func (g *G) Iter(val Val) (Val, E) {
  return val.Type(g).Iter(g, val)
}

func (g *G) Len(val Val) (Int, E) {
  return val.Type(g).Len(g, val)
}

func (g *G) Pop(val Val) (Val, Val, E) {
  return val.Type(g).Pop(g, val)
}

func (g *G) Print(val Val, out *strings.Builder) {
  val.Type(g).Print(g, val, out)
}

func (g *G) Push(val Val, its ...Val) (Val, E) {
  return val.Type(g).Push(g, val, its...)
}

func (g *G) Quote(task *Task, env *Env, val Val) (Val, E) {
  return val.Type(g).Quote(g, task, env, val)
}

func (g *G) Splat(val Val, out Vec) (Vec, E) {
  return val.Type(g).Splat(g, val, out)
}
