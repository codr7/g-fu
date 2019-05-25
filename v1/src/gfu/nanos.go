package gfu

import (
  "fmt"
  "bufio"
  //"log"
)

type Nanos Int

type NanosType struct {
  NumType
}

func (_ Nanos) Type(g *G) Type {
  return &g.NanosType
}

func (_ *NanosType) Div(g *G, x, y Val) (Val, E) {
  return g.IntType.Div(g, Int(x.(Nanos)), y)
}

func (_ *NanosType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "%vns", int64(val.(Nanos)))
  return nil
}

func (_ *NanosType) Float(g *G, val Val) (Float, E) {
  return g.IntType.Float(g, Int(val.(Nanos)))
}

func (_ *NanosType) Int(g *G, val Val) (Int, E) {
  return Int(val.(Nanos)), nil
}

func (_ *NanosType) Mul(g *G, x, y Val) (Val, E) {
  return g.IntType.Mul(g, Int(x.(Nanos)), y)
}

func (_ *NanosType) Sub(g *G, x, y Val) (Val, E) {
  yns, ok := y.(Nanos)

  if !ok {
    return nil, g.E("Expected Nanos: %v", y.Type(g))
  }

  out, e := g.IntType.Sub(g, Int(x.(Nanos)), Int(yns))

  if e != nil {
    return nil, e
  }
  
  return Nanos(out.(Int)), nil
}
