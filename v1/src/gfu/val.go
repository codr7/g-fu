package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Val interface {
  fmt.Stringer
  Dumper

  Bool(*G) bool
  Call(*G, *Task, *Env, Vec) (Val, E)
  Clone(*G) (Val, E)
  Drop(g *G, n Int) (Val, E)
  Dup(*G) (Val, E)
  Eq(*G, Val) bool
  Eval(*G, *Task, *Env) (Val, E)
  Expand(*G, *Task, *Env, Int) (Val, E)
  Extenv(*G, *Env, *Env, bool) E
  Is(*G, Val) bool
  Iter(*G) (Val, E)
  Len(*G) (Int, E)
  Pop(*G) (Val, Val, E)
  Print(*strings.Builder)
  Push(*G, ...Val) (Val, E)
  Quote(*G, *Task, *Env) (Val, E)
  Splat(*G, Vec) (Vec, E)
  Type(*G) *Type
}

type BasicVal struct {
  imp_type *Type
  imp      Val
}

func (v *BasicVal) Init(imp_type *Type, imp Val) *BasicVal {
  v.imp_type = imp_type
  v.imp = imp
  return v
}

func (_ BasicVal) Bool(g *G) bool {
  return true
}

func (v BasicVal) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E("Call not supported: ", v.imp_type)
}

func (v BasicVal) Clone(g *G) (Val, E) {
  return v.imp.Dup(g)
}

func (v BasicVal) Drop(g *G, n Int) (out Val, e E) {
  for i := Int(0); i < n; i++ {
    if _, out, e = v.imp.Pop(g); e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (v BasicVal) Dup(g *G) (Val, E) {
  return v.imp, nil
}

func (v BasicVal) Eq(g *G, rhs Val) bool {
  return v.imp.Is(g, rhs)
}

func (v BasicVal) Eval(g *G, task *Task, env *Env) (Val, E) {
  return v.imp, nil
}

func (v BasicVal) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  return v.imp, nil
}

func (v BasicVal) Extenv(g *G, src, dst *Env, clone bool) E {
  return nil
}

func (v BasicVal) Is(g *G, rhs Val) bool {
  return v.imp == rhs
}

func (v BasicVal) Iter(g *G) (Val, E) {
  return nil, g.E("Iter not supported: %v", v.imp_type)
}

func (v BasicVal) Len(g *G) (Int, E) {
  return -1, g.E("Len not supported: %v", v.imp_type)
}

func (v BasicVal) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: %v", v.imp_type)
}

func (v BasicVal) Print(out *strings.Builder) {
  v.imp.Dump(out)
}

func (v BasicVal) Push(g *G, its ...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", v.imp_type)
}

func (v BasicVal) Quote(g *G, task *Task, env *Env) (Val, E) {
  return v.imp, nil
}

func (v BasicVal) Splat(g *G, out Vec) (Vec, E) {
  return append(out, v.imp), nil
}

func (v BasicVal) Type(g *G) *Type {
  return v.imp_type
}

func (v BasicVal) String() string {
  return DumpString(v.imp)
}
