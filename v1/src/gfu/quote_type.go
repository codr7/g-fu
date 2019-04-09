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

func (t *QuoteType) Eval(g *G, pos Pos, val Val, env *Env) (v Val, e E) {
  return val.AsQuote().Quote(g, pos, env)
}

func (v Val) AsQuote() Val {
  return v.imp.(Val)
}
