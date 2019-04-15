package gfu

import (
  "fmt"
  //"log"
)

type Val interface {
  fmt.Stringer
  Dumper
  
  Bool(*G) bool
  Call(*G, *Task, *Env, Vec) (Val, E)
  Clone(*G) (Val, E)
  Dup(*G) (Val, E)
  Eq(*G, Val) bool
  Eval(*G, *Task, *Env) (Val, E)
  Is(*G, Val) bool
  Len(*G) (Int, E)
  Pop(*G) (Val, Val, E)
  Push(*G, ...Val) (Val, E)
  Quote(*G, *Task, *Env) (Val, E)
  Splat(*G, Vec) Vec
  Type(*G) *Type
}

type BasicVal struct {
  imp_type *Type
  imp Val
}

func (v *BasicVal) Init(imp_type *Type, imp Val) *BasicVal {
  v.imp_type = imp_type
  v.imp = imp
  return v
}

func (v BasicVal) Clone(g *G) (Val, E) {
  return v.imp.Dup(g)
}

func (v BasicVal) Dup(g *G) (Val, E) {
  return v.imp, nil
}

func (v BasicVal) Len(g *G) (Int, E) {
  return -1, g.E("Len not supported: %v", v.imp_type)
}

func (v *BasicVal) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: %v", v.imp_type)
}

func (v *BasicVal) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", v.imp_type)
}

func (v BasicVal) Type(g *G) *Type {
  return v.imp_type
}

func (v BasicVal) String() string {
  return DumpString(v.imp)
}

func (env *Env) AddVal(g *G, id string, val Val) {
  env.Let(g.Sym(id), val)
}
