package gfu

import (
	"bufio"
	//"log"
)

type Nil struct {
}

type NilType struct {
	BasicType
}

func (_ *Nil) Type(g *G) Type {
	return &g.NilType
}

func (_ *NilType) Bool(g *G, val Val) (bool, E) {
	return false, nil
}

func (_ *NilType) Dump(g *G, val Val, out *bufio.Writer) E {
	out.WriteRune('_')
	return nil
}
