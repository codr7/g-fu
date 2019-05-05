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
