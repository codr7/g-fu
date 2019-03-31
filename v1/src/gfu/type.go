package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  AsBool(g *G, val Val) bool
  Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error)
  Dump(x Val, out *strings.Builder)
  Eq(x, y Val) bool
  Splat(g *G, val Val, out []Val) []Val
}

type BasicType struct {
  id *Sym
}

func (t *BasicType) Init(id *Sym) *BasicType {
  t.id = id
  return t
}

func (t *BasicType) AsBool(g *G, val Val) bool {
  return true
}

func (t *BasicType) Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error) {
  if len(args) > 0 {
    return g.NIL, g.NewError(pos, "Too many args")
  }
  
  return val, nil
}

func (t *BasicType) Dump(x Val, out *strings.Builder) {
  fmt.Fprintf(out, "%v", x.imp)
}

func (t *BasicType) Eq(x, y Val) bool {
  return x.imp == y.imp
}

func (t *BasicType) Splat(g *G, val Val, out []Val) []Val {
  return append(out, val)
}
