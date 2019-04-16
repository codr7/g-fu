package gfu

import (
  //"log"
  "strings"
)

type Splat struct {
  Wrap
}

func NewSplat(g *G, val Val) (s Splat) {
  s.Wrap.Init(&g.SplatType, s, val)
  return s
}

func (s Splat) Dump(out *strings.Builder) {
  s.val.Dump(out)
  out.WriteString("..")
}

func (s Splat) Eq(g *G, rhs Val) bool {
  rs, ok := rhs.(Splat)

  if !ok {
    return false
  }

  return s.val.Eq(g, rs.val)
}

func (s Splat) Eval(g *G, task *Task, env *Env) (Val, E) {
  var e E
  s.val, e = s.val.Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  return s, nil
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
