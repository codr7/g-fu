package gfu

import (
  //"log"
  "strings"
)

type Splice struct {
  BasicWrap
}

type SpliceType struct {
  BasicWrapType
}

func NewSplice(g *G, val Val) (s Splice) {
  s.BasicWrap.Init(val)
  return s
}

func (_ Splice) Type(g *G) Type {
  return &g.SpliceType
}

func (_ *SpliceType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('%')
  return g.Dump(val.(Splice).val, out)
}

func (_ *SpliceType) Eval(g *G, task *Task, env *Env, val Val) (Val, E) {
  return nil, g.E("Unquoted splice")
}

func (_ *SpliceType) Quote(g *G, task *Task, env *Env, val Val) (Val, E) {
  s := val.(Splice)
  
  if v, ok := s.val.(Vec); ok {
    if len(v) == 1 {
      if sv, ok := v[0].(Splat); ok {
        var v Val
        var e E

        if v, e = g.Eval(task, env, sv.val); e != nil {
          return nil, e
        }

        return NewSplat(g, v), nil
      }
    }
  }

  return g.Eval(task, env, s.val)
}

func (_ *SpliceType) Unwrap(val Val) (Val, E) {
  return val.(Splice).val, nil
}

func (_ *SpliceType) Wrap(g *G, val Val) (Val, E) {
  return NewSplice(g, val), nil
}
