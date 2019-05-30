package gfu

import (
  "bufio"
)

type Throw struct {
  val Val
}

func (t Throw) Dump(g *G, out *bufio.Writer) (e E) {
  out.WriteString("Throw: ")

  if e = g.Dump(t.val, out); e != nil {
    return e
  }

  return nil
}

func (g *G) Throw(val Val) E {
  return Throw{val}
}
