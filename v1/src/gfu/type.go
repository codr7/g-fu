package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  AsBool(*G, Val) bool
  Call(*G, Pos, Val, ListForm, *Env) (Val, Error)
  Dump(Val, *strings.Builder)
  Eq(*G, Val, Val) bool
  Id() *Sym
  Is(*G, Val, Val) bool
  New(*G, Pos, Val, ListForm, *Env) (Val, Error)
  Splat(*G, Val, []Val) []Val
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

func (t *BasicType) Call(g *G, pos Pos, val Val, args ListForm, env *Env) (Val, Error) {
  if len(args) > 0 {
    return g.NIL, g.E(pos, "Too many args")
  }
  
  return val, nil
}

func (t *BasicType) Dump(x Val, out *strings.Builder) {
  fmt.Fprintf(out, "%v", x.imp)
}

func (t *BasicType) Eq(g *G, x Val, y Val) bool {
  return t.Is(g, x, y)
}

func (t *BasicType) Id() *Sym {
  return t.id
}

func (t *BasicType) Is(g *G, x Val, y Val) bool {
  return x == y
}

func (t *BasicType) New(g *G, pos Pos, al Val, args ListForm, env *Env) (Val, Error)  {
  return g.NIL, g.E(pos, "Missing constructor: %v", t.Id())
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
