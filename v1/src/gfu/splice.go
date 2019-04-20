package gfu

import (
  //"log"
  "strings"
)

type Splice struct {
  Wrap
}

func NewSplice(g *G, val Val) *Splice {
  s := new(Splice)
  s.Wrap.Init(&g.SpliceType, s, val)
  return s
}

func (s *Splice) Dump(out *strings.Builder) {
  out.WriteRune('%')
  s.val.Dump(out)
}

func (s *Splice) Eq(g *G, rhs Val) bool {
  rs, ok := rhs.(*Splice)

  if !ok {
    return false
  }

  return s.val.Eq(g, rs.val)
}

func (_ Splice) Eval(g *G, task *Task, env *Env) (Val, E) {
  return nil, g.E("Unquoted splice")
}

func (s *Splice) Is(g *G, rhs Val) bool {
  return s == rhs
}

func (s *Splice) Quote(g *G, task *Task, env *Env) (Val, E) {
  return s.val.Eval(g, task, env)
}

func (s *Splice) Splat(g *G, out Vec) Vec {
  return append(out, s)
}

func (_ Splice) Type(g *G) *Type {
  return &g.SpliceType
}
