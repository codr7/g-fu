package gfu

import (
  //"log"
  "strings"
)

type Splat struct {
  Wrap
}

func NewSplat(val Val) (s Splat) {
  s.val = val
  return s
}

func (s Splat) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return s, nil 
}

func (s Splat) Dump(out *strings.Builder) {
  s.val.Dump(out)
  out.WriteString("..")
}

func (s Splat) Eq(g *G, rhs Val) bool {
  return s.val.Is(g, rhs.(Splat).val)
}

func (s Splat) Eval(g *G, task *Task, env *Env) (Val, E) {
  var e E
  s.val, e = s.val.Eval(g, task, env)

  if e != nil {
    return nil, e
  }
  
  return s, nil
}

func (s Splat) Is(g *G, rhs Val) bool {
  return s == rhs
}

func (s Splat) Quote(g *G, task *Task, env *Env) (Val, E) {
  var e E
  s.val, e = s.val.Quote(g, task, env)

  if e != nil {
    return nil, e
  }

  return s, nil
}

func (s Splat) Splat(g *G, out Vec) Vec {
  v := s.val

  if _, ok := v.(Vec); !ok {
    return append(out, s)
  }

  return v.Splat(g, out)
}

func (s Splat) Type(g *G) *Type {
  return &g.SplatType
}

