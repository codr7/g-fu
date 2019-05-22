package gfu

import (
  "fmt"
  "bufio"
  //"log"
)

type NSecs Int

type NSecsType struct {
  IntType
}

func (_ NSecs) Type(g *G) Type {
  return &g.NSecsType
}

func (_ *NSecsType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "%vns", int64(val.(NSecs)))
  return nil
}

func (t *NSecsType) Sub(g *G, x, y Val) (Val, E) {
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
