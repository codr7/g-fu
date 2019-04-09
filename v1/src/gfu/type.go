package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  Bool(*G, Val) bool
  Call(*G, Pos, Val, []Val, *Env) (Val, E)
  Dump(Val, *strings.Builder)
  Eq(*G, Val, Val) bool
  Eval(*G, Pos, Val, *Env) (Val, E)
  Id() *Sym
  Init(*Sym)
  Is(*G, Val, Val) bool
  New(*G, Pos, Val, []Val, *Env) (Val, E)
  Quote(*G, Pos, Val, *Env) (Val, E)
  Splat(*G, Pos, Val, []Val) []Val
}

type BasicType struct {
  id *Sym
}

func (t *BasicType) Init(id *Sym) {
  t.id = id
}

func (t *BasicType) Bool(g *G, val Val) bool {
  return true
}

func (t *BasicType) Call(g *G, pos Pos, val Val, args []Val, env *Env) (Val, E) {
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

func (t *BasicType) Eval(g *G, pos Pos, val Val, env *Env) (Val, E) {
  return val, nil
}

func (t *BasicType) Id() *Sym {
  return t.id
}

func (t *BasicType) Is(g *G, x Val, y Val) bool {
  return x.imp == y.imp
}

func (t *BasicType) New(g *G, pos Pos, val Val, args []Val, env *Env) (Val, E)  {
  return g.NIL, g.E(pos, "Missing constructor: %v", t.Id())
}

func (t *BasicType) Quote(g *G, pos Pos, val Val, env *Env) (Val, E) {
  return val, nil
}

func (t *BasicType) Splat(g *G, pos Pos, val Val, out []Val) []Val {
  return append(out, val)
}

func (e *Env) AddType(g *G, id string, t Type) Type {
  t.Init(g.Sym(id))
  var v Val
  var mt Type

  if id == "Meta" {
    mt = t
  } else {
    mt = g.MetaType
  }
  
  v.Init(NIL_POS, mt, t)
  e.Put(t.Id(), v)
  return t
}
