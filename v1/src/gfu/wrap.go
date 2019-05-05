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
  Unwrap(Val) Val
  Wrap(*G, Val) Val
}

type BasicWrapType struct {
  BasicType
}

func (g *G) Unwrap(val Val) Val {
  t := val.Type(g)
  wt, ok := t.(WrapType)

  if !ok {
    return nil
  }
  
  return wt.Unwrap(val)
}

func (g *G) Wrap(typ Type, val Val) (Val, E) {
  wt, ok := typ.(WrapType)

  if !ok {
    return nil, g.E("Wrap not supported: %v", typ)
  }
  
  return wt.Wrap(g, val), nil
}

func (w *BasicWrap) Init(val Val) *BasicWrap {
  w.val = val
  return w
}

func (_ *BasicWrapType) Bool(g *G, val Val) (bool, E) {
  v := g.Unwrap(val)

  if v == nil {
    return false, g.E("Unwrap not supported: %v", val.Type(g))
  }
  
  return g.Bool(v)
}

func (_ *BasicWrapType) Clone(g *G, val Val) (v Val, e E) {
  v = g.Unwrap(val)
  
  if v == nil {
    return nil, g.E("Unwrap not supported: %v", val.Type(g))
  }

  if v, e = g.Clone(v); e != nil {
    return nil, e
  }

  return g.Wrap(val.Type(g), v)
}

func (_ *BasicWrapType) Dump(g *G, val Val, out *strings.Builder) E {
  v := g.Unwrap(val)

  if v == nil {
    return g.E("Unwrap not supported: %v", val.Type(g))
  }

  return g.Dump(v, out)
}

func (_ *BasicWrapType) Dup(g *G, val Val) (v Val, e E) {
  v = g.Unwrap(val)
  
  if v == nil {
    return nil, g.E("Unwrap not supported: %v", val.Type(g))
  }

  if v, e = g.Dup(v); e != nil {
    return nil, e
  }

  return g.Wrap(val.Type(g), v)
}

func (_ *BasicWrapType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lv := g.Unwrap(lhs)

  if lv == nil {
    return false, nil
  }  

  rv := g.Unwrap(rhs)

  if rv == nil {
    return false, nil
  }  

  return g.Eq(lv, rv)
}

func (_ *BasicWrapType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  v := g.Unwrap(val)

  if v == nil {
    return g.E("Unwrap not supported: %v", val.Type(g))
  }

  return g.Extenv(src, dst, v, clone)
}
