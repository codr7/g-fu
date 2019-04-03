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

func (t *VecType) Eq(g *G, x Val, y Val) bool {
  xv, yv := x.AsVec().items, y.AsVec().items

  if len(xv) != len(yv) {
    return false
  }

  for i, xi := range xv {
    yi := yv[i]
    
    if !xi.Eq(g, yi) {
      return false
    }
  }

  return true
}

func (t *VecType) New(g *G, pos Pos, val Val, args VecForm, env *Env) (Val, Error)  {
  is, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, g.E(pos, "Constructor arg eval failed: %v", e)
  }

  var out Val
  v := new(Vec)
  v.items = is
  out.Init(g.Vec, v)
  return out, nil
}

func (t *VecType) Splat(g *G, val Val, out []Val) []Val {
  return append(out, val.AsVec().items...)
}

func (v Val) AsVec() *Vec {
  return v.imp.(*Vec)
}

