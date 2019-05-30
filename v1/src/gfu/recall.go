package gfu

import (
  "bufio"
)

type Recall struct {
  args Vec
}

func NewRecall(args Vec) (r Recall) {
  r.args = args
  return r
}

func (r Recall) Dump(g *G, out *bufio.Writer) (e E) {
  out.WriteString("(recall")

  for _, a := range r.args {
    if e = g.Dump(a, out); e != nil {
      return e
    }
  }

  out.WriteRune(')')
  return nil
}

func (r Recall) String() string {
  return "Recall"
}
