package gfu

import (
  //"log"
  "strings"
)

type False struct {
  BasicVal
}

func (f *False) Init(g *G) *False {
  f.BasicVal.Init(&g.FalseType, f)
  return f
}

func (_ *False) Bool(g *G) bool {
  return false
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
