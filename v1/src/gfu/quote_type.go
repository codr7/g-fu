package gfu

import (
  //"log"
  "strings"
)

type QuoteType struct {
  BasicType
}

func (t *QuoteType) Bool(g *G, val Val) bool {
  return val.AsQuote().AsBool(g)
}

func (t *QuoteType) Dump(val Val, out *strings.Builder) {
  out.WriteRune('\'')
  val.AsQuote().Dump(out)
}

func (t *QuoteType) Eq(g *G, x Val, y Val) bool {
  return x.AsQuote().Is(g, y.AsQuote())
}

func (t *QuoteType) Eval(g *G, val Val, env *Env) (Val, E) {
  var e E
  
  if val, e = val.AsQuote().Quote(g, env); e != nil {
    return g.NIL, e
  }

  if val.val_type == g.VecType {
    val.imp = val.Splat(g, nil)
  }

  return val, nil
}

func (v Val) AsQuote() Val {
  return v.imp.(Val)
}
