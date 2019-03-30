package gfu

import (
  "strings"
)

type VecType struct {
  BasicType
}

func (t *VecType) Init(id *Sym) *VecType {
  t.BasicType.Init(id)
  return t
}

func (t *VecType) AsBool(g *G, val Val) bool {
  return val.AsVec().items != nil
}

func (t *VecType) Dump(val Val, out *strings.Builder) {
  v := val.AsVec()
  out.WriteRune('(')
  
  for i, iv := range v.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    iv.Dump(out)
  }
  
  out.WriteRune(')')
}

func (v Val) AsVec() *Vec {
  return v.imp.(*Vec)
}

