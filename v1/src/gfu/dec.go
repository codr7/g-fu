package gfu

import (
  "fmt"
  //"log"
  "math/big"
  "strings"
)

type Dec big.Float

type DecType struct {
  NumType
}

func (d *Dec) Abs() {
  f := big.Float(*d)
  f.Abs(&f)
  *d = Dec(f)
}

func (d *Dec) Add(val Dec) {
  x, y := big.Float(*d), big.Float(val)
  x.Add(&x, &y)
  *d = Dec(x)
}

func (d Dec) Cmp(val Dec) int {
  x, y := big.Float(d), big.Float(val)
  return x.Cmp(&y)
}

func (d *Dec) Div(val Dec) {
  x, y := big.Float(*d), big.Float(val)
  x.Quo(&x, &y)
  *d = Dec(x)
}

func (d Dec) Int() Int {
  f := big.Float(d)
  i, _ := f.Int64()
  return Int(i)
}

func (d *Dec) Mul(val Dec) {
  x, y := big.Float(*d), big.Float(val)
  x.Mul(&x, &y)
  *d = Dec(x)
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

func (d *Dec) SetFloat(val float64) {
  var f big.Float
  f.SetFloat64(val)
  *d = Dec(f)
}

func (d *Dec) SetInt(val Int) {
  var f big.Float
  f.SetInt64(int64(val))
  *d = Dec(f)
}

func (d Dec) Sign() int {
  f := big.Float(d)
  return f.Sign()
}

func (d Dec) String() string {
  f := big.Float(d)
  return f.String()
}

func (d *Dec) Sub(val Dec) {
  x, y := big.Float(*d), big.Float(val)
  x.Sub(&x, &y)
  *d = Dec(x)
}

func (_ Dec) Type(g *G) Type {
  return &g.DecType
}

func (t *DecType) Abs(g *G, x Val) (Val, E) {
  xd := x.(Dec)
  xd.Abs()
  return xd, nil
}

func (t *DecType) Add(g *G, x, y Val) (Val, E) {
  xd := x.(Dec)
  yd, ok := y.(Dec)

  if !ok {
    return nil, g.E("Expected Dec: ", y.Type(g))
  }

  xd.Add(yd)
  return xd, nil
}

func (_ *DecType) Bool(g *G, val Val) (bool, E) {
  return val.(Dec).Sign() != 0, nil
}

func (_ *DecType) Dec(g *G, val Val) (Dec, E) {
  return val.(Dec), nil
}

func (t *DecType) Div(g *G, x, y Val) (Val, E) {
  xd := x.(Dec)
  yd, ok := y.(Dec)

  if !ok {
    return nil, g.E("Expected Dec: ", y.Type(g))
  }

  xd.Div(yd)
  return xd, nil
}

func (_ *DecType) Dump(g *G, val Val, out *strings.Builder) E {
  fmt.Fprintf(out, "%v", val.(Dec))
  return nil
}

func (_ *DecType) Int(g *G, val Val) (Int, E) {
  return val.(Dec).Int(), nil
}

func (t *DecType) Mul(g *G, x, y Val) (Val, E) {
  xd := x.(Dec)
  yd, ok := y.(Dec)

  if !ok {
    return nil, g.E("Expected Dec: ", y.Type(g))
  }

  xd.Mul(yd)
  return xd, nil
}

func (t *DecType) Neg(g *G, x Val) (Val, E) {
  xd := x.(Dec)
  xd.Neg()
  return x, nil
}

func (t *DecType) Is(g *G, x, y Val) bool {
  xd := x.(Dec)
  yd, ok := y.(Dec)

  return ok && xd.Cmp(yd) == 0
}

func (t *DecType) Sub(g *G, x, y Val) (Val, E) {
  xd := x.(Dec)
  yd, ok := y.(Dec)

  if !ok {
    return nil, g.E("Expected Dec: ", y.Type(g))
  }

  xd.Sub(yd)
  return xd, nil
}
