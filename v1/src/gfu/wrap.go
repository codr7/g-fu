package gfu

import (
  //"log"
  "strings"
)

type Wrap struct {
  val Val
}

func (w Wrap) Bool(g *G) bool {
  return w.val.Bool(g)
}

func (w Wrap) Dump(out *strings.Builder) {
  w.val.Dump(out)
}

func (w Wrap) String() string {
  return DumpString(w.val)
}
