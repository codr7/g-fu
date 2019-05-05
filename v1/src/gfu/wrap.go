package gfu

import (
  //"log"
  "strings"
)

type BasicWrap struct {
  val Val
}

type WrapType interface {
  Type
  Unwrap(Val) (Val, E)
  Wrap(*G, Val) (Val, E)
}

type BasicWrapType struct {
  BasicType
}

func (g *G) Unwrap(val Val) (Val, E) {
  t := val.Type(g)
  wt, ok := t.(WrapType)

  if !ok {
    return nil, g.E("Unwrap not supported: %v", t)
  }
  
  return wt.Unwrap(val)
}

func (g *G) Wrap(typ Type, val Val) (Val, E) {
  wt, ok := typ.(WrapType)

  if !ok {
    return nil, g.E("Wrap not supported: %v", typ)
  }
  
  return wt.Wrap(g, val)
}

func (w *BasicWrap) Init(val Val) *BasicWrap {
  w.val = val
  return w
}

func (_ *BasicWrapType) Bool(g *G, val Val) (bool, E) {
  v, e := g.Unwrap(val)

  if e != nil {
    return false, e
  }
  
  return g.Bool(v)
}

func (_ *BasicWrapType) Clone(g *G, val Val) (Val, E) {
  v, e := g.Unwrap(val)

  if e != nil {
    return nil, e
  }

  if v, e = g.Clone(v); e != nil {
    return nil, e
  }

  return g.Wrap(val.Type(g), v)
}

func (_ *BasicWrapType) Dump(g *G, val Val, out *strings.Builder) E {
  v, e := g.Unwrap(val)

  if e != nil {
    return e
  }  

  return g.Dump(v, out)
}

func (_ *BasicWrapType) Dup(g *G, val Val) (Val, E) {
  v, e := g.Unwrap(val)

  if e != nil {
    return nil, e
  }  

  if v, e = g.Dup(v); e != nil {
    return nil, e
  }

  return g.Wrap(val.Type(g), v)
}

func (_ *BasicWrapType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lv, e := g.Unwrap(lhs)

  if e != nil {
    return false, e
  }  

  rv, e := g.Unwrap(rhs)

  if e != nil {
    return false, e
  }  

  return g.Eq(lv, rv)
}

func (_ *BasicWrapType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  v, e := g.Unwrap(val)

  if e != nil {
    return e
  }  

  return g.Extenv(src, dst, v, clone)
}
