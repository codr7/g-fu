package gfu

import (
  //"log"
  "strings"
)

type VecType struct {
  BasicType
}

func (t *VecType) Bool(g *G, val Val) bool {
  return len(val.AsVec()) > 0
}

func (t *VecType) Dump(val Val, out *strings.Builder) {
  v := val.AsVec()
  out.WriteRune('(')
  
  for i, iv := range v {
    if i > 0 {
      out.WriteRune(' ')
    }

    iv.Dump(out)
  }
  
  out.WriteRune(')')
}

func (t *VecType) Eq(g *G, x Val, y Val) bool {
  xv, yv := x.AsVec(), y.AsVec()

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

func (t *VecType) Eval(g *G, val Val, env *Env) (Val, E) {
  v := val.AsVec()
  
  if len(v) == 0 {
    return g.NIL, nil
  }
  
  first := v[0]
  first_val, e := first.Eval(g, env)
  
  if e != nil {
    return g.NIL, e
  }

  result, e := first_val.Call(g, v[1:], env)
  
  if e != nil {
    return g.NIL, g.E("Call failed: %v", e)
  }
  
  return result, nil
}

func (t *VecType) Is(g *G, x Val, y Val) bool {
  return t.Eq(g, x, y)
}

func (t *VecType) Quote(g *G, val Val, env *Env) (Val, E) {
  var out Vec

  for _, v := range val.AsVec() {
    qv, e := v.Quote(g, env)

    if e != nil {
      return g.NIL, e
    }
    
    out = out.Push(qv)
  }

  var v Val
  v.Init(g.VecType, out)
  return v, nil
}

func (t *VecType) Splat(g *G, val Val, out Vec) Vec {
  for _, it := range val.AsVec() {
    if it.val_type == g.SplatType {
      out = it.Splat(g, out)
    } else {
      if it.val_type == g.VecType {
        it.imp = it.Splat(g, nil)
      }
      
      out = append(out, it)
    }
  }

  return out
}

func (v Val) AsVec() Vec {
  return v.imp.(Vec)
}

