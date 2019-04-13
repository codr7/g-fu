package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Int int

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
  
func (i Int) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "%v", i)
}

func (i Int) Eq(g *G, rhs Val) bool {
  ri, ok := rhs.(Int)
  return ok && ri == i
}

func (i Int) Eval(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Is(g *G, rhs Val) bool {
  return i.Eq(g, rhs)
}

func (i Int) Quote(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Splat(g *G, out Vec) Vec {
  return append(out, i)
}

func (_ Int) Type(g *G) *Type {
  return &g.IntType
}
