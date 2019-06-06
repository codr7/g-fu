package gfu

import (
	"bufio"
	"fmt"
	//"log"
	"math/big"
)

type Float big.Float

type FloatType struct {
	NumType
}

func (f *Float) Abs() {
	bf := big.Float(*f)
	bf.Abs(&bf)
	*f = Float(bf)
}

func (f *Float) Add(val Float) {
	x, y := big.Float(*f), big.Float(val)
	x.Add(&x, &y)
	*f = Float(x)
}

func (f Float) Cmp(val Float) int {
	x, y := big.Float(f), big.Float(val)
	return x.Cmp(&y)
}

func (f *Float) Div(val Float) {
	x, y := big.Float(*f), big.Float(val)
	x.Quo(&x, &y)
	*f = Float(x)
}

func (f *Float) Mul(val Float) {
	x, y := big.Float(*f), big.Float(val)
	x.Mul(&x, &y)
	*f = Float(x)
}

func (f *Float) Neg() {
	bf := big.Float(*f)
	bf.Neg(&bf)
	*f = Float(bf)
}

func (f *Float) Parse(g *G, in string) E {
	var bf big.Float

	if _, _, e := bf.Parse(in, 10); e != nil {
		return g.E("Failed parsing Float: %v", e)
	}

	*f = Float(bf)
	return nil
}

func (f *Float) SetFloat(val float64) {
	var bf big.Float
	bf.SetFloat64(val)
	*f = Float(bf)
}

func (f *Float) SetInt(val Int) {
	var bf big.Float
	bf.SetInt64(int64(val))
	*f = Float(bf)
}

func (f Float) Sign() int {
	bf := big.Float(f)
	return bf.Sign()
}

func (f Float) String() string {
	bf := big.Float(f)
	return bf.String()
}

func (f *Float) Sub(val Float) {
	x, y := big.Float(*f), big.Float(val)
	x.Sub(&x, &y)
	*f = Float(x)
}

func (_ Float) Type(g *G) Type {
	return &g.FloatType
}

func (t *FloatType) Abs(g *G, x Val) (Val, E) {
	xf := x.(Float)
	xf.Abs()
	return xf, nil
}

func (t *FloatType) Add(g *G, x, y Val) (Val, E) {
	xf := x.(Float)

	switch y := y.(type) {
	case Float:
		xf.Add(y)
	case Int:
		var yf Float
		yf.SetInt(y)
		xf.Add(yf)
	default:
		yf, e := g.Float(y)

		if e != nil {
			return nil, e
		}

		xf.Add(yf)
	}

	return xf, nil
}

func (_ *FloatType) Bool(g *G, val Val) (bool, E) {
	return val.(Float).Sign() != 0, nil
}

func (_ *FloatType) Float(g *G, val Val) (Float, E) {
	return val.(Float), nil
}

func (t *FloatType) Div(g *G, x, y Val) (Val, E) {
	xf := x.(Float)

	switch y := y.(type) {
	case Float:
		xf.Div(y)
	case Int:
		var yf Float
		yf.SetInt(y)
		xf.Div(yf)
	default:
		yf, e := g.Float(y)

		if e != nil {
			return nil, e
		}

		xf.Div(yf)
	}

	return xf, nil
}

func (_ *FloatType) Dump(g *G, val Val, out *bufio.Writer) E {
	fmt.Fprintf(out, "%v", val.(Float))
	return nil
}

func (t *FloatType) Mul(g *G, x, y Val) (Val, E) {
	xf := x.(Float)

	switch y := y.(type) {
	case Float:
		xf.Mul(y)
	case Int:
		var yf Float
		yf.SetInt(y)
		xf.Mul(yf)
	default:
		yf, e := g.Float(y)

		if e != nil {
			return nil, e
		}

		xf.Mul(yf)
	}

	return xf, nil
}

func (t *FloatType) Neg(g *G, x Val) (Val, E) {
	xf := x.(Float)
	xf.Neg()
	return xf, nil
}

func (t *FloatType) Is(g *G, x, y Val) bool {
	xf := x.(Float)
	yf, ok := y.(Float)

	return ok && xf.Cmp(yf) == 0
}

func (t *FloatType) Sub(g *G, x, y Val) (Val, E) {
	xf := x.(Float)
	yf, ok := y.(Float)

	if !ok {
		return nil, g.E("Expected Float: %v", y.Type(g))
	}

	xf.Sub(yf)
	return xf, nil
}
