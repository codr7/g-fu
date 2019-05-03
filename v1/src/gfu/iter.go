package gfu

import (
	"strings"
)

type BasicIter struct {
	BasicVal
}

func (i *BasicIter) Dump(out *strings.Builder) {
	out.WriteString(i.imp_type.id.name)
}

func (i *BasicIter) Iter(g *G) (out Val, e E) {
	if out, e = i.imp.Clone(g); e != nil {
		return nil, e
	}

	return out, nil
}
