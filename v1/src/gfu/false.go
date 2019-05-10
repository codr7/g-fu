package gfu

import (
	//"log"
	"strings"
)

type False struct {
}

type FalseType struct {
	BasicType
}

func (_ *False) Type(g *G) Type {
	return &g.FalseType
}

func (_ *FalseType) Bool(g *G, val Val) (bool, E) {
	return false, nil
}

func (_ *FalseType) Dump(g *G, val Val, out *strings.Builder) E {
	out.WriteRune('F')
	return nil
}
