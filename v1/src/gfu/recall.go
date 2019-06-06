package gfu

import (
	"bufio"
)

type Recall struct {
	args Vec
}

type RecallType struct {
	BasicType
}

func NewRecall(args Vec) (r Recall) {
	r.args = args
	return r
}

func (_ Recall) Type(g *G) Type {
	return &g.RecallType
}

func (_ RecallType) Dump(g *G, val Val, out *bufio.Writer) (e E) {
	r := val.(Recall)
	out.WriteString("(recall")

	for _, a := range r.args {
		if e = g.Dump(a, out); e != nil {
			return e
		}
	}

	out.WriteRune(')')
	return nil
}
