package gfu

import (
  //"log"
  "strings"
)

type Quote struct {
  BasicWrap
}

type QuoteType struct {
  BasicWrapType
}

func NewQuote(g *G, val Val) *Quote {
  q := new(Quote)
  q.BasicWrap.Init(val)
  return q
}

func (_ *Quote) Type(g *G) Type {
  return &g.QuoteType
}

func (_ *QuoteType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('\'')
  return g.Dump(val.(*Quote).val, out)
}

func (_ *QuoteType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lq := lhs.(*Quote)
  rq, ok := rhs.(*Quote)

  if !ok {
    return false, nil
  }
  
  return g.Eq(lq.val, rq.val)
}

func (_ *QuoteType) Eval(g *G, task *Task, env *Env, val Val) (Val, E) {
  q := val.(*Quote)
  qv, e := g.Quote(task, env, q.val)

  if e != nil {
    return nil, e
  }

  if v, ok := qv.(Vec); ok {
    if qv, e = g.Splat(v, nil); e != nil {
      return nil, e
    }
  }

  return qv, nil
}

func (_ *QuoteType) Quote(g *G, task *Task, env *Env, val Val) (Val, E) {
  q := val.(*Quote)
  
  if _, ok := q.val.(*Splice); !ok {
    return q, nil
  }

  var v Val
  var e E

  if v, e = g.Quote(task, env, q.val); e != nil {
    return nil, e
  }

  return NewQuote(g, v), nil
}

func (_ *QuoteType) Unwrap(val Val) (*BasicWrap, E) {
  return &val.(*Quote).BasicWrap, nil
}
