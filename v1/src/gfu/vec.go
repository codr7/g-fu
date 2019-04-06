package gfu

import (
  //"log"
)

type Vec struct {
  items []Val
}

func (v *Vec) Push(its...Val) {
  v.items = append(v.items, its...)
}

func (v *Vec) Peek(g *G) Val {
  is := v.items
  n := len(is)
  
  if n == 0 {
    return g.NIL
  }

  return is[n-1]
}

func (v *Vec) Pop(g *G) Val {
  is := v.items
  n := len(is)

  if n == 0 {
    return g.NIL
  }

  var it Val
  it, v.items = is[n-1], is[:n-1]
  return it
}
