package gfu

import (
  "strings"
)

type Val struct {
  val_type Type
  imp interface{}
}

func NewVal(val_type Type, imp interface{}) *Val {
  return new(Val).Init(val_type, imp)
}

func (v *Val) Init(val_type Type, imp interface{}) *Val {
  v.val_type = val_type
  v.imp = imp
  return v
}

func (v Val) Call(g *G, args []Val, env *Env, pos Pos) (Val, Error) {
  return v.val_type.Call(g, v, args, env, pos)
}

func (v Val) Dump(out *strings.Builder) {
  v.val_type.Dump(v, out)
}

func (v Val) String() string {
  var out strings.Builder
  v.Dump(&out)
  return out.String()
}
