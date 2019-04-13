package gfu

import (
  //"log"
  "strings"
)

type Vec []Val

func (v Vec) Bool(g *G) bool {
  return len(v) > 0
}

func (v Vec) Call(g *G, args Vec, env *Env) (Val, E) {
  return v, nil 
}

func (v Vec) Dump(out *strings.Builder) {
  out.WriteRune('(')
  
  for i, iv := range v {
    if i > 0 {
      out.WriteRune(' ')
    }

    iv.Dump(out)
  }
  
  out.WriteRune(')')
}

func (v Vec) Eq(g *G, rhs Val) bool {
  rv := rhs.(Vec)
  
  if len(v) != len(rv) {
    return false
  }

  for i, vi := range v {
    ri := rv[i]
    
    if !vi.Is(g, ri) {
      return false
    }
  }

  return true
}

func (v Vec) Eval(g *G, env *Env) (Val, E) {
  if len(v) == 0 {
    return g.NIL, nil
  }
  
  first := v[0]
  first_val, e := first.Eval(g, env)
  
  if e != nil {
    return nil, e
  }

  result, e := first_val.Call(g, v[1:], env)
  
  if e != nil {
    return nil, g.E("Call failed: %v", e)
  }
  
  return result, nil
}

func (v Vec) EvalExpr(g *G, env *Env) (Val, E) {
  var out Val = g.NIL
  
  for _, it := range v {
    var e E
    
    if out, e = it.Eval(g, env); e != nil {
      return nil, e
    }

    if g.recall {
      break
    }
  }

  return out, nil
}

func (v Vec) EvalVec(g *G, env *Env) (Vec, E) {
  var out Vec
  
  for _, it := range v {
    it, e := it.Eval(g, env)

    if e != nil {
      return nil, g.E("Arg eval failed: %v", e)
    }

    if g.recall {
      break
    }
    
    if _, ok := it.(Splat); ok {
      out = it.Splat(g, out)
    } else {
      if _, ok := it.(Vec); ok {
        it = it.Splat(g, nil)
      }
      
      out = append(out, it)
    }
  }

  return out, nil
}

func (v Vec) Is(g *G, rhs Val) bool {
  return v.Eq(g, rhs)
}

func (v Vec) Len() Int {
  return Int(len(v))
}

func (v Vec) Push(its...Val) Vec {
  return append(v, its...)
}

func (v Vec) Peek(g *G) Val {
  n := len(v)
  
  if n == 0 {
    return g.NIL
  }

  return v[n-1]
}

func (v Vec) Pop(g *G) (Val, Vec) {
  n := len(v)

  if n == 0 {
    return g.NIL, v
  }

  return v[n-1], v[:n-1]
}

func (v Vec) Quote(g *G, env *Env) (Val, E) {
  var out Vec

  for _, it := range v {
    q, e := it.Quote(g, env)

    if e != nil {
      return nil, e
    }
    
    out = out.Push(q)
  }

  return out, nil
}

func (v Vec) Splat(g *G, out Vec) Vec {
  for i, it := range v {
    if _, ok := it.(Splat); ok {
      out = it.Splat(g, out)
    } else {
      if _, ok := it.(Vec); ok {
        it = it.Splat(g, nil)
        v[i] = it
      }
      
      out = append(out, it)
    }
  }

  return out
}

func (v Vec) Type(g *G) *Type {
  return &g.VecType
}
