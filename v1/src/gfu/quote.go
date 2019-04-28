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
    if qv, e = v.Splat(g, nil); e != nil {
      return nil, e
    }
  }

  return qv, nil
}

func (q *Quote) Is(g *G, rhs Val) bool {
  return q == rhs
}

func (q *Quote) Quote(g *G, task *Task, env *Env) (Val, E) {
  if _, ok := q.val.(*Splice); !ok {
    return q, nil
  }

  var v Val
  var e E

  if v, e = q.val.Quote(g, task, env); e != nil {
    return nil, e
  }
  
  return NewQuote(g, v), nil
}

func (_ Quote) Type(g *G) *Type {
  return &g.QuoteType
}
