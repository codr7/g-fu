package gfu

import (
  //"log"
  "strings"
)

type Nil struct {
}

type NilType struct {
  BasicType
}

type NilIter struct {
}

var nil_iter NilIter

type NilIterType struct {
  BasicIterType
}

func (_ *Nil) Type(g *G) Type {
  return &g.NilType
}

func (_ *NilType) Bool(g *G, val Val) (bool, E) {
  return false, nil
}

func (_ *NilType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('_')
  return nil
}

func (_ *NilType) Iter(g *G, val Val) (Val, E) {
  return &nil_iter, nil
}

func (_ *NilType) Len(g *G, val Val) (Int, E) {
  return 0, nil
}

func (_ *NilType) Pop(g *G, val Val) (Val, Val, E) {
  return val, val, nil
}

func (_ *NilType) Push(g *G, val Val, its ...Val) (Val, E) {
  return Vec(its), nil
}

func (_ *NilType) Splat(g *G, val Val, out Vec) (Vec, E) {
  return out, nil
}

func (_ *NilIter) Type(g *G) Type {
  return &g.NilIterType
}

func (_ *NilIterType) Bool(g *G, val Val) (bool, E) {
  return false, nil
}

func (_ *NilIterType) Drop(g *G, val Val, n Int) (Val, E) {
  return nil, g.E("Nothing to drop")
}

func (_ *NilIterType) Dup(g *G, val Val) (Val, E) {
  out := *val.(*NilIter)
  return &out, nil
}

func (_ *NilIterType) Eq(g *G, lhs, rhs Val) (bool, E) {
  _, ok := rhs.(*NilIter)
  return ok, nil
}

func (_ *NilIterType) Pop(g *G, val Val) (Val, Val, E) {
  return &g.NIL, val, nil
}

func (_ *NilIterType) Splat(g *G, val Val, out Vec) (Vec, E) {
  return out, nil
}
