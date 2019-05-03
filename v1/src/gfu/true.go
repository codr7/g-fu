package gfu

import (
	//"log"
	"strings"
)

type True struct {
	BasicVal
}

func (t *True) Init(g *G) *True {
	t.BasicVal.Init(&g.TrueType, t)
	return t
}

func (_ *True) Dump(out *strings.Builder) {
	out.WriteRune('T')
}
