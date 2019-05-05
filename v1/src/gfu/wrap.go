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
  Unwrap(Val) (*BasicWrap, E)
}

type BasicWrapType struct {
  BasicType
}

func (g *G) Unwrap(val Val) (*BasicWrap, E) {
  t := val.Type(g)
  wt, ok := t.(WrapType)

  if !ok {
    return nil, g.E("Unwrap not supported: %v", t)
  }
  
  return wt.Unwrap(val)
}

func (w *BasicWrap) Init(val Val) *BasicWrap {
  w.val = val
  return w
}

func (_ *BasicWrapType) Bool(g *G, val Val) (bool, E) {
  w, e := g.Unwrap(val)

  if e != nil {
    return false, e
  }
  
  return g.Bool(w.val)
}

func (_ *BasicWrapType) Clone(g *G, val Val) (Val, E) {
  w, e := g.Unwrap(val)

  if e != nil {
    return nil, e
  }

  if w.val, e = g.Clone(w.val); e != nil {
    return nil, e
  }

  return val, nil
}

func (_ *BasicWrapType) Dump(g *G, val Val, out *strings.Builder) E {
  w, e := g.Unwrap(val)

  if e != nil {
    return e
  }  

  return g.Dump(w.val, out)
}

func (_ *BasicWrapType) Dup(g *G, val Val) (Val, E) {
  w, e := g.Unwrap(val)

  if e != nil {
    return nil, e
  }  

  if w.val, e = g.Dup(w.val); e != nil {
    return nil, e
  }

  return val, nil
}

func (_ *BasicWrapType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lw, e := g.Unwrap(lhs)

  if e != nil {
    return false, e
  }  

  rw, e := g.Unwrap(rhs)

  if e != nil {
    return false, e
  }  

  return g.Eq(lw.val, rw.val)
}

func (_ *BasicWrapType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  w, e := g.Unwrap(val)

  if e != nil {
    return e
  }  

  return g.Extenv(src, dst, w.val, clone)
}
