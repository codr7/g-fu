package gfu

import ()

type BasicIterType struct {
  BasicType
}

type IterType struct {
  BasicIterType
}

func (_ *BasicIterType) Iter(g *G, val Val) (out Val, e E) {
  return g.Clone(val)
}
