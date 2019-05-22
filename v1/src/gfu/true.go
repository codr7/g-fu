package gfu

import (
  "bufio"
  //"log"
)

type True struct {
}

type TrueType struct {
  BasicType
}

func (_ *True) Type(g *G) Type {
  return &g.TrueType
}

func (_ *TrueType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteRune('T')
  return nil
}
