package gfu

import (
  //"log"
  "strings"
)

type Vec []Val

func (v Vec) Bool(g *G) bool {
  return len(v) > 0
}

func (v Vec) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E("Call not supported: Vec")
}

func (v Vec) Clone(g *G) (Val, E) {
  out := make(Vec, len(v))
  var e E
  
  for i, it := range v {
    if out[i], e = it.Clone(g); e != nil {
      return nil, e
    }
  }
  
  return out, nil
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

func (v Vec) Dup(g *G) (Val, E) {
  out := make(Vec, len(v))
  copy(out, v)
  return out, nil
}

func (v Vec) Eq(g *G, rhs Val) bool {
  rv, ok := rhs.(Vec)

  if !ok || len(v) != len(rv) {
    return false
  }

  for i, vi := range v {
    ri := rv[i]

    if !vi.Eq(g, ri) {
      return false
    }
  }

  return true
}

func (v Vec) Eval(g *G, task *Task, env *Env) (Val, E) {
  if len(v) == 0 {
    return &g.NIL, nil
  }

  fv, e := v[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  result, e := fv.Call(g, task, env, v[1:])

  if e != nil {
    return nil, e
  }

  return result, nil
}

func (v Vec) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  if len(v) == 0 {
    return &g.NIL, nil
  }

  idv := v[0]

  if s, ok := idv.(*Sym); ok && s == g.nil_sym {
    return &g.NIL, nil
  }

  id, ok := idv.(*Sym)

  if !ok {
    return v, nil
  }

  _, mv := env.Find(id)
  
  if mv == nil {
    return v, v.ExpandVec(g, task, env, depth-1)
  }
  
  m, ok := mv.Val.(*Mac)

  if !ok {
    return v, v.ExpandVec(g, task, env, depth-1)
  }
  
  out, e := m.ExpandCall(g, task, env, v[1:])

  if depth == 1 || e != nil {
    return out, e
  }

  return out.Expand(g, task, env, depth-1)
}

func (v Vec) EvalExpr(g *G, task *Task, env *Env) (Val, E) {
  var out Val = &g.NIL

  for _, it := range v {
    var e E

    if out, e = it.Eval(g, task, env); e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (v Vec) EvalVec(g *G, task *Task, env *Env) (Vec, E) {
  var out Vec
  var e E
  
  for _, it := range v {
    it, e = it.Eval(g, task, env)

    if e != nil {
      return nil, e
    }

    if _, ok := it.(*Splat); ok {
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

func (v Vec) ExpandVec(g *G, task *Task, env *Env, depth Int) E {
  for i, it := range v {
    var e E

    if v[i], e = it.Expand(g, task, env, depth); e != nil {
      return e
    }
  }

  return nil
}

func (v Vec) Extenv(g *G, src, dst *Env, clone bool) E {
  for _, it := range v {
    if e := it.Extenv(g, src, dst, clone); e != nil {
      return e
    }
  }

  return nil
}

func (v Vec) Is(g *G, rhs Val) bool {
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

func (v Vec) Len(g *G) (Int, E) {
  return Int(len(v)), nil
}

func (v Vec) Push(g *G, its...Val) (Val, E) {
  return append(v, its...), nil
}

func (v Vec) Peek(g *G) Val {
  n := len(v)

  if n == 0 {
    return &g.NIL
  }

  return v[n-1]
}

func (v Vec) Pop(g *G) (Val, Val, E) {
  n := len(v)

  if n == 0 {
    return &g.NIL, v, nil
  }

  return v[n-1], v[:n-1], nil
}

func (v Vec) Print(out *strings.Builder) {
  for i, iv := range v {
    if i > 0 {
      out.WriteRune(' ')
    }
    
    iv.Print(out)
  }
}

func (v Vec) Quote(g *G, task *Task, env *Env) (Val, E) {
  var e E
  out := make(Vec, len(v))

  for i, it := range v {
    out[i], e = it.Quote(g, task, env)

    if e != nil {
      return nil, e
    }
  }

  return out, nil
}

func(v Vec) Reverse() Vec {
  for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
    v[i], v[j] = v[j], v[i]
  }

  return v
}

func (v Vec) Splat(g *G, out Vec) Vec {  
  for _, it := range v {
    if _, ok := it.(*Splat); ok {
      out = it.Splat(g, out)
    } else {
      if _, ok := it.(Vec); ok {
        it = it.Splat(g, nil)
      }

      out = append(out, it)
    }
  }

  return out
}

func (v Vec) String() string {
  return DumpString(v)
}

func (v Vec) Type(g *G) *Type {
  return &g.VecType
}
