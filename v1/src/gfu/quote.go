package gfu

import (
  //"log"
  "strings"
)

type Quote struct {
  Wrap
}

func NewQuote(g *G, val Val) *Quote {
  q := new(Quote)
  q.Wrap.Init(&g.QuoteType, q, val)
  return q
}

func (q *Quote) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return q, nil
}

func (q *Quote) Dump(out *strings.Builder) {
  out.WriteRune('\'')
  q.val.Dump(out)
}

func (q *Quote) Eq(g *G, rhs Val) bool {
  rq, ok := rhs.(*Quote)

  if !ok {
    return false
  }

  return q.val.Eq(g, rq.val)
}

func (q *Quote) Eval(g *G, task *Task, env *Env) (Val, E) {
  qv, e := q.val.Quote(g, task, env)

  if e != nil {
    return nil, e
  }

  if v, ok := qv.(Vec); ok {
    qv = v.Splat(g, nil)
  }

  return qv, nil
}

func (q *Quote) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  return q, nil
}

func (q *Quote) Is(g *G, rhs Val) bool {
  return q == rhs
}

func (q *Quote) Quote(g *G, task *Task, env *Env) (Val, E) {
  return q, nil
}

func (q *Quote) Splat(g *G, out Vec) Vec {
  return append(out, q)
}

func (_ Quote) Type(g *G) *Type {
  return &g.QuoteType
}
