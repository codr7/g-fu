package gfu

import (
	"bufio"
	"fmt"
	//"log"
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
		v := Int(xb) + y

		if v >= 0 && v < 256 {
			return Byte(v), nil
		}

		return v, nil
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

func (_ *ByteType) Dump(g *G, val Val, out *bufio.Writer) E {
	fmt.Fprintf(out, "0x%02x", val.(Byte))
	return nil
}

func (_ *ByteType) Eq(g *G, lhs, rhs Val) (bool, E) {
	lb := lhs.(Byte)
	rb, e := g.Byte(rhs)

	if e != nil {
		return false, e
	}

	return lb == rb, nil
}

func (_ *ByteType) Int(g *G, val Val) (Int, E) {
	return Int(val.(Byte)), nil
}

func (_ *ByteType) Print(g *G, val Val, out *bufio.Writer) E {
	out.WriteByte(byte(val.(Byte)))
	return nil
}

func (_ *ByteType) Sub(g *G, x Val, y Val) (Val, E) {
	xb := x.(Byte)

	switch y := y.(type) {
	case Byte:
		return xb - y, nil
	case Int:
		v := Int(xb) - y

		if v >= 0 && v < 256 {
			return Byte(v), nil
		}

		return v, nil
	default:
		break
	}

	return nil, g.E("Invalid sub arg: %v", y.Type(g))
}
