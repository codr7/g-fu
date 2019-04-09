package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type PrimType struct {
  BasicType
}

func (t *PrimType) Call(g *G, pos Pos, val Val, args []Val, env *Env) (Val, E) {
  p := val.AsPrim()
  
  if e := p.arg_list.Check(g, pos, args); e != nil {
    return g.NIL, e
  }

  return p.imp(g, pos, args, env)
}

func (t *PrimType) Dump(val Val, out *strings.Builder) {
  fmt.Fprintf(out, "(Prim %v)", val.AsPrim().id)
}

func (v Val) AsPrim() *Prim {
  return v.imp.(*Prim)
}
