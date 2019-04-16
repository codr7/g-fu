package gfu

import (
  //"log"
  "strings"
)

type Move struct {
  Wrap
}

func NewMove(g *G, val Val) (m Move) {
  m.Wrap.Init(&g.MoveType, m, val)
  return m
}

func (m Move) Clone(g *G) (Val, E) {
  return m.Dup(g)
}

func (s Move) Dump(out *strings.Builder) {
  out.WriteString("(move ")
  s.val.Dump(out)
  out.WriteRune(')')
}

func (m Move) Dup(g *G) (Val, E) {
  if m.val == nil {
    return nil, g.E("Dup move")
  }
  
  var v Val
  v, m.val = m.val, nil
  return v, nil
}
