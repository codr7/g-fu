package gfu

import (
	//"log"
	"strings"
)

type True struct {
}

type TrueType struct {
	BasicType
}

func (_ *True) Type(g *G) Type {
	return &g.TrueType
}

func (_ *TrueType) Dump(g *G, val Val, out *strings.Builder) E {
	out.WriteRune('T')
	return nil
}
