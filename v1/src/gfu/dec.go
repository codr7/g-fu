package gfu

import (
  "fmt"
  //"log"
  "math/big"
  "strings"
)

type Dec big.Float

type DecType struct {
  BasicType
}

func (d *Dec) Neg() {
  f := big.Float(*d)
  f.Neg(&f)
  *d = Dec(f)
}

func (d *Dec) Parse(g *G, in string) E {
  var f big.Float

  if _, _, e := f.Parse(in, 10); e != nil {
    return g.E("Failed parsing Dec: %v", e)
  }
  
  *d = Dec(f)
  return nil
}

func (d Dec) Sign() int {
  f := big.Float(d)
  return f.Sign()
}

func (d Dec) String() string {
  f := big.Float(d)
  return f.String()
}

func (_ Dec) Type(g *G) Type {
  return &g.DecType
}

func (_ *DecType) Bool(g *G, val Val) (bool, E) {
  return val.(Dec).Sign() != 0, nil
}

func (_ *DecType) Dump(g *G, val Val, out *strings.Builder) E {
  fmt.Fprintf(out, "%v", val.(Dec))
  return nil
}
