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
  Id() *Sym
  New(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error)
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

func (t *BasicType) Id() *Sym {
  return t.id
}

func (t *BasicType) New(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error)  {
  return g.NIL, g.NewError(pos, "Missing constructor: %v", t.Id())
}

func (t *BasicType) Splat(g *G, val Val, out []Val) []Val {
  return append(out, val)
}

func (e *Env) AddType(g *G, t Type) Type {
  var v Val
  v.Init(g.Meta, t)
  e.Put(t.Id(), v)
  return t
}
