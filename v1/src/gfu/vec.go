package gfu

type Vec struct {
  items []Val
}

func (v *Vec) Push(item Val) {
  v.items = append(v.items, item)
}

func (v *Vec) Pop(g *G) Val {
  if v.items == nil {
    return g.NIL
  }

  is := v.items
  n := len(is)
  var it Val
  it, v.items = is[n-1], is[:n-1]
  return it
}
