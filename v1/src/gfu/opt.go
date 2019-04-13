package gfu

import (
  //"log"
  "strings"
)

type Opt struct {
  Wrap
}

func NewOpt(val Val) (o Opt) {
  o.val = val
  return o
}

func (o Opt) Call(g *G, args Vec, env *Env) (Val, E) {
  return o, nil 
}

func (o Opt) Dump(out *strings.Builder) {
  o.val.Dump(out)
  out.WriteRune('?')
}

func (o Opt) Eq(g *G, rhs Val) bool {
  return o.val.Is(g, rhs.(Opt).val)
}

func (o Opt) Eval(g *G, env *Env) (Val, E) {
  v, e := o.val.Eval(g, env)

  if e != nil {
    return nil, e
  }
  
  return NewOpt(v), nil
}

func (o Opt) Is(g *G, rhs Val) bool {
  return o == rhs
}

func (o Opt) Quote(g *G, env *Env) (Val, E) {
  v, e := o.val.Quote(g, env)

  if e != nil {
    return nil, e
  }

  return NewOpt(v), nil
}

func (o Opt) Splat(g *G, out Vec) Vec {
  return append(out, o)
}

func (_ Opt) Type(g *G) *Type {
  return &g.OptType
}
