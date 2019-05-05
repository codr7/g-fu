package gfu

import (
)

type BasicIterType struct {
  BasicType
}

func (_ *BasicIterType) Iter(g *G, val Val) (out Val, e E) {
  return g.Clone(val)
}
