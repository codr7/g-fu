package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type PrimType struct {
  BasicType
}

func (t *PrimType) Call(g *G, pos Pos, val Val, args []Form, env *Env) (Val, E) {
  pp := g.prim
  p := val.AsPrim()
  g.prim = p
  v, e := p.imp(g, pos, args, env)
  g.prim = pp
  return v, e
}

func (t *PrimType) Dump(val Val, out *strings.Builder) {
  fmt.Fprintf(out, "(Prim %v)", val.AsPrim().id)
}

func (v Val) AsPrim() *Prim {
  return v.imp.(*Prim)
}
