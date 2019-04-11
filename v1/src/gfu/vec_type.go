package gfu

import (
  //"log"
  "strings"
)

type VecType struct {
  BasicType
}

func (t *VecType) Bool(g *G, val Val) bool {
  return len(val.AsVec().items) > 0
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
    
    if !xi.Is(g, yi) {
      return false
    }
  }

  return true
}

func (t *VecType) Eval(g *G, pos Pos, val Val, env *Env) (Val, E) {
  v := val.AsVec()
  
  if len(v.items) == 0 {
    return g.NIL, nil
  }
  
  first := v.items[0]
  first_val, e := first.Eval(g, pos, env)
  
  if e != nil {
    return g.NIL, e
  }

  result, e := first_val.Call(g, pos, v.items[1:], env)
  
  if e != nil {
    return g.NIL, g.E(pos, "Call failed: %v", e)
  }
  
  return result, nil
}

func (t *VecType) New(g *G, pos Pos, val Val, args []Val, env *Env) (Val, E)  {
  var out Val
  v := new(Vec)
  v.items = args
  out.Init(pos, g.VecType, v)
  return out, nil
}

func (t *VecType) Quote(g *G, pos Pos, val Val, env *Env) (Val, E) {
  var out Vec

  for _, v := range val.AsVec().items {
    qv, e := v.Quote(g, pos, env)

    if e != nil {
      return g.NIL, e
    }
    
    out.Push(qv)
  }

  var v Val
  v.Init(pos, g.VecType, &out)
  return v, nil
}

func (t *VecType) Splat(g *G, pos Pos, val Val, out []Val) []Val {
  for _, it := range val.AsVec().items {
    if it.val_type == g.SplatType {
      out = it.Splat(g, pos, out)
    } else {
      out = append(out, it)
    }
  }

  return out
}

func (v Val) AsVec() *Vec {
  return v.imp.(*Vec)
}

