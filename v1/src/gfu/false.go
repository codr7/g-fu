package gfu

import (
	//"log"
	"strings"
)

type False struct {
	BasicVal
}

func (f *False) Init(g *G) *False {
	f.BasicVal.Init(&g.FalseType, f)
	return f
}

func (_ *False) Bool(g *G) bool {
	return false
}

func (_ *False) Dump(out *strings.Builder) {
	out.WriteRune('F')
}
