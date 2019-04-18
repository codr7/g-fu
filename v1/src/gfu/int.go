package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Int int64

func (i Int) Abs() Int {
  if i < 0 {
    return -i
  }

  return i
}

func (i Int) Bool(g *G) bool {
  return i != 0
}

func (i Int) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return i, nil
}

func (i Int) Clone(g *G) (Val, E) {
  return i, nil
}

func (i Int) Dup(g *G) (Val, E) {
  return i, nil
}

func (i Int) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "%v", int64(i))
}

func (i Int) Eq(g *G, rhs Val) bool {
  return i.Is(g, rhs)
}

func (i Int) Eval(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  return i, nil
}

func (i Int) Is(g *G, rhs Val) bool {
  return i == rhs
}

func (_ Int) Len(g *G) (Int, E) {
  return -1, g.E("Len not supported: Int")
}

func (_ Int) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: Int")
}

func (_ Int) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: Int")
}

func (i Int) Quote(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Splat(g *G, out Vec) Vec {
  return append(out, i)
}

func (i Int) String() string {
  return DumpString(i)
}

func (_ Int) Type(g *G) *Type {
  return &g.IntType
}
