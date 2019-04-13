package gfu

import (
  //"log"
  "strings"
)

type False struct {
}

func (_ *False) Bool(g *G) bool {
  return false
}

func (f *False) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return f, nil
}
  
func (_ *False) Dump(out *strings.Builder) {
  out.WriteRune('F')
}

func (f *False) Eq(g *G, rhs Val) bool {
  return f == rhs
}

func (f *False) Eval(g *G, task *Task, kenv *Env) (Val, E) {
  return f, nil
}

func (f *False) Is(g *G, rhs Val) bool {
  return f == rhs
}

func (f *False) Quote(g *G, task *Task, env *Env) (Val, E) {
  return f, nil
}

func (f *False) Splat(g *G, out Vec) Vec {
  return append(out, f)
}

func (f *False) String() string {
  return "F"
}

func (_ *False) Type(g *G) *Type {
  return &g.FalseType
}
