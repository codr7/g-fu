package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Byte byte

type ByteType struct {
  NumType
}

func (_ Byte) Type(g *G) Type {
  return &g.ByteType
}

func (_ *ByteType) Add(g *G, x Val, y Val) (Val, E) {
  xb := x.(Byte)
  yb, e := g.Byte(y)
  
  if e != nil {
    return nil, e
  }

  return xb+yb, nil
}

func (_ *ByteType) Bool(g *G, val Val) (bool, E) {
  return val.(Byte) != 0, nil
}

func (_ *ByteType) Byte(g *G, val Val) (Byte, E) {
  return val.(Byte), nil
}

func (_ *ByteType) Dump(g *G, val Val, out *strings.Builder) E {
  fmt.Fprintf(out, "0x%02x", val.(Byte))
  return nil
}
