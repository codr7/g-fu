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

  switch y := y.(type) {
  case Byte:
    return xb + y, nil    
  case Int:
    return Byte(Int(xb) + y), nil    
  default:
    break
  }
  
  return nil, g.E("Invalid add arg: %v", y.Type(g))
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

func (_ *ByteType) Int(g *G, val Val) (Int, E) {
  return Int(val.(Byte)), nil
}

func (_ *ByteType) Sub(g *G, x Val, y Val) (Val, E) {
  xb := x.(Byte)

  switch y := y.(type) {
  case Byte:
    return xb - y, nil    
  case Int:
    return Byte(Int(xb) - y), nil    
  default:
    break
  }
  
  return nil, g.E("Invalid sub arg: %v", y.Type(g))
}
