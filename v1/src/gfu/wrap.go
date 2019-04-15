package gfu

import (
  //"log"
  "strings"
)

type Wrap struct {
  BasicVal
  val Val
}

func (w *Wrap) Init(imp_type *Type, imp Val, val Val) *Wrap {
  w.BasicVal.Init(imp_type, imp)
  w.val = val
  return w
}
  
func (w Wrap) Bool(g *G) bool {
  return w.val.Bool(g)
}

func (w Wrap) Dump(out *strings.Builder) {
  w.val.Dump(out)
}

func (w Wrap) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", w.imp_type)
}

func (w Wrap) String() string {
  return DumpString(w.imp)
}

func (w Wrap) Type(g *G) *Type {
  return w.imp_type
}
