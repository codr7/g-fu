package gfu

import (
  "fmt"
  "strings"
)

type Type interface {
  Dump(x interface{}, out *strings.Builder)
  Eq(x, y interface{}) bool
}

type BasicType struct {
  id *Sym
}

func (t *BasicType) Init(id *Sym) *BasicType {
  t.id = id
  return t
}

func (t *BasicType) Dump(x interface{}, out *strings.Builder) {
  fmt.Fprintf(out, "%v", x)
}

func (t *BasicType) Eq(x, y interface{}) bool {
  return x == y
}
