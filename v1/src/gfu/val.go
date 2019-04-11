package gfu

import (
  //"log"
  "strings"
)

type Val struct {
  pos Pos
  val_type Type
  imp interface{}
}

func (v *Val) Init(pos Pos, val_type Type, imp interface{}) *Val {
  v.pos = pos
  v.val_type = val_type
  v.imp = imp
  return v
}
  
func (v Val) Call(g *G, pos Pos, args Vec, env *Env) (Val, E) {
  return v.val_type.Call(g, pos, v, args, env)
}

func (v Val) Dump(out *strings.Builder) {
  v.val_type.Dump(v, out)
}

func (v Val) Eq(g *G, rhs Val) bool {
  return v.val_type.Eq(g, v, rhs)
}

func (v Val) Eval(g *G, pos Pos, env *Env) (Val, E) {
  return v.val_type.Eval(g, pos, v, env)
}

func (v Val) Is(g *G, rhs Val) bool {
  return v.val_type.Is(g, v, rhs)
}

func (v Val) Quote(g *G, pos Pos, env *Env) (Val, E) {
  return v.val_type.Quote(g, pos, v, env)
}

func (v Val) Splat(g *G, pos Pos, out Vec) Vec {
  return v.val_type.Splat(g, pos, v, out)
}

func (v Val) String() string {
  return DumpString(v)
}

func (env *Env) AddVal(g *G, id string, val_type Type, val interface{}, out *Val) {
  out.Init(NIL_POS, val_type, val)
  env.Put(g.Sym(id), *out)
}
