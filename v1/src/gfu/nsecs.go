package gfu

import (
  "fmt"
  "bufio"
  //"log"
)

type NSecs Int

type NSecsType struct {
  NumType
}

func (_ NSecs) Type(g *G) Type {
  return &g.NSecsType
}

func (_ *NSecsType) Div(g *G, x, y Val) (Val, E) {
  return g.IntType.Div(g, Int(x.(NSecs)), y)
}

func (_ *NSecsType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "%vns", int64(val.(NSecs)))
  return nil
}

func (_ *NSecsType) Float(g *G, val Val) (Float, E) {
  return g.IntType.Float(g, Int(val.(NSecs)))
}

func (_ *NSecsType) Int(g *G, val Val) (Int, E) {
  return Int(val.(NSecs)), nil
}

func (_ *NSecsType) Mul(g *G, x, y Val) (Val, E) {
  return g.IntType.Mul(g, Int(x.(NSecs)), y)
}

func (_ *NSecsType) Sub(g *G, x, y Val) (Val, E) {
  yns, ok := y.(NSecs)

  if !ok {
    return nil, g.E("Expected NSecs: %v", y.Type(g))
  }

  out, e := g.IntType.Sub(g, Int(x.(NSecs)), Int(yns))

  if e != nil {
    return nil, e
  }
  
  return NSecs(out.(Int)), nil
}
