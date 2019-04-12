package gfu

import (
  //"log"
  "strings"
)

type Val struct {
  val_type Type
  imp interface{}
}

func (v *Val) Init(val_type Type, imp interface{}) *Val {
  v.val_type = val_type
  v.imp = imp
  return v
}
  
func (v Val) Call(g *G, args Vec, env *Env) (Val, E) {
  return v.val_type.Call(g, v, args, env)
}

func (v Val) Dump(out *strings.Builder) {
  v.val_type.Dump(v, out)
}

func (v Val) Eq(g *G, rhs Val) bool {
  return v.val_type.Eq(g, v, rhs)
}

func (v Val) Eval(g *G, env *Env) (Val, E) {
  return v.val_type.Eval(g, v, env)
}

func (v Val) Is(g *G, rhs Val) bool {
  return v.val_type.Is(g, v, rhs)
}

func (v Val) Quote(g *G, env *Env) (Val, E) {
  return v.val_type.Quote(g, v, env)
}

func (v Val) Splat(g *G, out Vec) Vec {
  return v.val_type.Splat(g, v, out)
}

func (v Val) String() string {
  return DumpString(v)
}

func (env *Env) AddVal(g *G, id string, val_type Type, val interface{}, out *Val) {
  out.Init(val_type, val)
  env.Put(g.Sym(id), *out)
}
