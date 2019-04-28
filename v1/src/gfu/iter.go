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

func (i *BasicIter) Iter(g *G) (Val, E) {
  return i.imp, nil
}

