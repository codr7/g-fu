package gfu

import (
  //"log"
  "strings"
)

type True struct {
}

func (_ *True) Bool(g *G) bool {
  return true
}

func (t *True) Call(g *G, args Vec, ent *Env) (Val, E) {
  return t, nil
}
  
func (_ *True) Dump(out *strings.Builder) {
  out.WriteRune('T')
}

func (t *True) Eq(g *G, rhs Val) bool {
  return t == rhs
}

func (t *True) Eval(g *G, ent *Env) (Val, E) {
  return t, nil
}

func (t *True) Is(g *G, rhs Val) bool {
  return t == rhs
}

func (t *True) Quote(g *G, env *Env) (Val, E) {
  return t, nil
}

func (t *True) Splat(g *G, out Vec) Vec {
  return append(out, t)
}

func (t *True) String() string {
  return "T"
}

func (_ *True) Type(g *G) *Type {
  return &g.TrueType
}
