package gfu

import (
  //"log"
  "strings"
)

type Vec []Val

type VecType struct {
  BasicType
}

type VecIter struct {
  in  Vec
  pos Int
}

type VecIterType struct {
  BasicIterType
}

func (v Vec) Delete(i int) Vec {
  return append(v[:i], v[i+1:]...)
}

func (v Vec) EvalExpr(g *G, task *Task, env *Env) (Val, E) {
  var out Val = &g.NIL

  for _, it := range v {
    var e E

    if out, e = g.Eval(task, env, it); e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (v Vec) EvalVec(g *G, task *Task, env *Env) (Vec, E) {
  var out Vec
  var e E

  for _, it := range v {
    it, e = g.Eval(task, env, it)

    if e != nil {
      return nil, e
    }

    if _, ok := it.(Splat); ok {
      if out, e = g.Splat(it, out); e != nil {
        return nil, e
      }
    } else {
      if _, ok := it.(Vec); ok {
        if it, e = g.Splat(it, nil); e != nil {
          return nil, e
        }
      }

      out = append(out, it)
    }
  }

  return out, nil
}

func (v Vec) ExpandVec(g *G, task *Task, env *Env, depth Int) (Vec, E) {
  for i, it := range v {
    var e E

    if v[i], e = g.Expand(task, env, it, depth); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func (v Vec) Len() Int {
  return Int(len(v))
}

func (v Vec) Peek(g *G) Val {
  n := len(v)

  if n == 0 {
    return &g.NIL
  }

  return v[n-1]
}

func (v Vec) PopKey(g *G, key *Sym) (Val, Val, E) {
  for i := 0; i < len(v); {
    k, ok := v[i].(*Sym)

    if !ok {
      return nil, nil, g.E("Invalid key: %v", v[i].Type(g))
    }

    if k == key {
      v = v.Delete(i)
      val := v[i]
      v = v.Delete(i)
      return val, v, nil
    } else {
      i += 2
    }
  }

  return &g.NIL, v, nil
}

func (v Vec) Reverse() Vec {
  for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
    v[i], v[j] = v[j], v[i]
  }

  return v
}

func (_ Vec) Type(g *G) Type {
  return &g.VecType
}

func (_ *VecType) Bool(g *G, val Val) (bool, E) {
  return val.(Vec).Len() > 0, nil
}

func (_ *VecType) Clone(g *G, val Val) (Val, E) {
  v := val.(Vec)
  dst := make(Vec, len(v))
  var e E

  for i, it := range v {
    if dst[i], e = g.Clone(it); e != nil {
      return nil, e
    }
  }

  return dst, nil
}

func (_ *VecType) Drop(g *G, val Val, n Int) (Val, E) {
  v := val.(Vec)
  vl := Int(len(v))

  if vl < n {
    return nil, g.E("Nothing to drop")
  }

  return v[:vl-n], nil
}

func (_ *VecType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('(')

  for i, iv := range val.(Vec) {
    if i > 0 {
      out.WriteRune(' ')
    }

    if e := g.Dump(iv, out); e != nil {
      return e
    }
  }

  out.WriteRune(')')
  return nil
}

func (_ *VecType) Dup(g *G, val Val) (Val, E) {
  v := val.(Vec)
  out := make(Vec, len(v))
  copy(out, v)
  return out, nil
}

func (_ *VecType) Eq(g *G, lhs, rhs Val) (bool, E) {
  lv := lhs.(Vec)
  rv, ok := rhs.(Vec)

  if !ok || len(lv) != len(rv) {
    return false, nil
  }

  for i, li := range lv {
    ri := rv[i]

    if ok, e := g.Eq(li, ri); e != nil || !ok {
      return ok, e
    }
  }

  return true, nil
}

func (_ *VecType) Eval(g *G, task *Task, env *Env, val Val) (Val, E) {
  v := val.(Vec)

  if len(v) == 0 {
    return Vec(nil), nil
  }

  target, ce := v[0], env

  if id, ok := target.(*Sym); ok {
    if id == g.nop_sym {
      return &g.NIL, nil
    }

    var e E

    if target, ce, e = id.Lookup(g, task, env, false); e != nil {
      return nil, e
    }
  }

  return g.Call(task, ce, target, v[1:], env)
}

func (_ *VecType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (Val, E) {
  v := val.(Vec)
  n := len(v)

  if n == 0 {
    return Vec(nil), nil
  }

  idv := v[0]
  id, ok := idv.(*Sym)

  if !ok {
    return v, nil
  }

  if id == g.nop_sym {
    return val, nil
  }

  if id == g.Sym("do") && n < 3 {
    if n == 1 {
      return &g.NIL, nil
    }

    return g.Expand(task, env, v[1], depth)
  }

  _, mv := env.Find(id)

  if mv == nil {
    return v.ExpandVec(g, task, env, depth-1)
  }

  m, ok := mv.Val.(*Mac)

  if !ok {
    return v.ExpandVec(g, task, env, depth-1)
  }

  out, e := m.ExpandCall(g, task, env, v[1:])

  if depth == 1 || e != nil {
    return out, e
  }

  return g.Expand(task, env, out, depth-1)
}

func (_ *VecType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  v := val.(Vec)
  
  if len(v) > 1 && v[0] == g.set_sym {
    if k, ok := v[1].(Vec); ok {
      if e := g.Extenv(src, dst, g.Sym("set-%v", k[0]), clone); e != nil {
        return e
      }
    }
  }
  
  for _, it := range v {
    if e := g.Extenv(src, dst, it, clone); e != nil {
      return e
    }
  }

  return nil
}

func (_ *VecType) Is(g *G, lhs, rhs Val) bool {
  lv := lhs.(Vec)
  rv, ok := rhs.(Vec)

  if !ok || len(lv) != len(rv) {
    return false
  }

  for i, li := range lv {
    ri := rv[i]

    if !g.Is(li, ri) {
      return false
    }
  }

  return true
}

func (_ *VecType) Iter(g *G, val Val) (Val, E) {
  return new(VecIter).Init(g, val.(Vec)), nil
}

func (_ *VecType) Len(g *G, val Val) (Int, E) {
  return val.(Vec).Len(), nil
}

func (_ *VecType) Push(g *G, val Val, its ...Val) (Val, E) {
  return append(val.(Vec), its...), nil
}

func (_ *VecType) Pop(g *G, val Val) (Val, Val, E) {
  v := val.(Vec)
  n := len(v)

  if n == 0 {
    return &g.NIL, v, nil
  }

  return v[n-1], v[:n-1], nil
}

func (_ *VecType) Print(g *G, val Val, out *strings.Builder) {
  for i, iv := range val.(Vec) {
    if i > 0 {
      out.WriteRune(' ')
    }

    g.Print(iv, out)
  }
}

func (_ *VecType) Quote(g *G, task *Task, env *Env, val Val) (Val, E) {
  v := val.(Vec)
  out := make(Vec, len(v))
  var e E

  for i, it := range v {
    out[i], e = g.Quote(task, env, it)

    if e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (_ *VecType) Splat(g *G, val Val, out Vec) (Vec, E) {
  var e E

  for _, it := range val.(Vec) {
    if _, ok := it.(Splat); ok {
      if out, e = g.Splat(it, out); e != nil {
        return nil, e
      }
    } else {
      if _, ok := it.(Vec); ok {
        if it, e = g.Splat(it, nil); e != nil {
          return nil, e
        }
      }

      out = append(out, it)
    }
  }

  return out, nil
}

func (i *VecIter) Init(g *G, in Vec) *VecIter {
  i.in = in
  return i
}

func (_ *VecIter) Type(g *G) Type {
  return &g.VecIterType
}

func (_ *VecIterType) Bool(g *G, val Val) (bool, E) {
  i := val.(*VecIter)
  return i.pos < Int(len(i.in)), nil
}

func (_ *VecIterType) Drop(g *G, val Val, n Int) (Val, E) {
  i := val.(*VecIter)

  if Int(len(i.in))-i.pos < n {
    return nil, g.E("Nothing to drop")
  }

  i.pos += n
  return i, nil
}

func (_ *VecIterType) Dup(g *G, val Val) (Val, E) {
  out := *val.(*VecIter)
  return &out, nil
}

func (_ *VecIterType) Eq(g *G, lhs, rhs Val) (bool, E) {
  li := lhs.(*VecIter)
  ri, ok := rhs.(*VecIter)

  if !ok {
    return false, nil
  }

  ok, e := g.Eq(ri.in, li.in)

  if e != nil {
    return false, e
  }

  return ok && ri.pos == li.pos, nil
}

func (_ *VecIterType) Pop(g *G, val Val) (Val, Val, E) {
  i := val.(*VecIter)

  if i.pos >= Int(len(i.in)) {
    return &g.NIL, i, nil
  }

  v := i.in[i.pos]
  i.pos++
  return v, i, nil
}

func (_ *VecIterType) Splat(g *G, val Val, out Vec) (Vec, E) {
  i := val.(*VecIter)
  return append(out, i.in[i.pos:]...), nil
}
