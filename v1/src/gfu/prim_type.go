package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type PrimType struct {
  BasicType
}

func (t *PrimType) Init(id *Sym) *PrimType {
  t.BasicType.Init(id)
  return t
}

func (t *PrimType) Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error) {
  p := val.AsPrim()
  nargs := len(args)
  
  if (p.min_args != -1 && nargs < p.min_args) ||
    (p.max_args != -1 && nargs > p.max_args) {
    return g.NIL, g.NewError(pos, "Arg mismatch")
  }

  return p.imp(g, args, env, pos)
}

func (t *PrimType) Dump(val Val, out *strings.Builder) {
  p := val.AsPrim()
  fmt.Fprintf(out, "(prim %v)", p.id)
}

func (v Val) AsPrim() *Prim {
  return v.imp.(*Prim)
}
