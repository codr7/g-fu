package gfu

import (
  //"log"
  "strings"
)

type Vec []Val

type VecIter struct {
  BasicIter
  in  Vec
  pos Int
}

func (v Vec) Bool(g *G) bool {
  return len(v) > 0
}

func (v Vec) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E("Call not supported: Vec")
}

func (v Vec) Clone(g *G) (Val, E) {
  dst := make(Vec, len(v))
  var e E

  for i, it := range v {
    if dst[i], e = it.Clone(g); e != nil {
      return nil, e
    }
  }

  return dst, nil
}

func (v Vec) Delete(i int) Vec {
  return append(v[:i], v[i+1:]...)
}

func (v Vec) Drop(g *G, n Int) (Val, E) {
  vl := Int(len(v))

  if vl < n {
    return nil, g.E("Nothing to drop")
  }

  return v[:vl-n], nil
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

  fid, ok := v[0].(*Sym)

  if !ok {
    return nil, g.E("Invalid call target: %v", v[0])
  }

  f, e := env.Get(g, fid)

  if e != nil {
    return nil, e
  }

  result, e := f.Call(g, task, env, v[1:])

  if e != nil {
    return nil, e
  }

  return result, nil
}

func (v Vec) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  n := len(v)

  if n == 0 {
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

  if id == g.Sym("do") && n < 3 {
    if n == 1 {
      return &g.NIL, nil
    }

    return v[1].Expand(g, task, env, depth)
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
      if out, e = it.Splat(g, out); e != nil {
        return nil, e
      }
    } else {
      if _, ok := it.(Vec); ok {
        if it, e = it.Splat(g, nil); e != nil {
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

    if v[i], e = it.Expand(g, task, env, depth); e != nil {
      return nil, e
    }
  }

  return v, nil
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

func (v Vec) Iter(g *G) (Val, E) {
  return new(VecIter).Init(g, v), nil
}

func (v Vec) Len(g *G) (Int, E) {
  return Int(len(v)), nil
}

func (v Vec) Push(g *G, its ...Val) (Val, E) {
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

func (v Vec) Reverse() Vec {
  for i, j := 0, len(v)-1; i < j; i, j = i+1, j-1 {
    v[i], v[j] = v[j], v[i]
  }

  return v
}

func (v Vec) Splat(g *G, out Vec) (Vec, E) {
  var e E

  for _, it := range v {
    if _, ok := it.(*Splat); ok {
      if out, e = it.Splat(g, out); e != nil {
        return nil, e
      }
    } else {
      if _, ok := it.(Vec); ok {
        if it, e = it.Splat(g, nil); e != nil {
          return nil, e
        }
      }

      out = append(out, it)
    }
  }

  return out, nil
}

func (v Vec) String() string {
  return DumpString(v)
}

func (v Vec) Type(g *G) *Type {
  return &g.VecType
}

func (i *VecIter) Init(g *G, in Vec) *VecIter {
  i.BasicVal.Init(&g.IterType, i)
  i.in = in
  return i
}

func (i *VecIter) Bool(g *G) bool {
  return i.pos < Int(len(i.in))
}

func (i *VecIter) Drop(g *G, n Int) (Val, E) {
  if Int(len(i.in))-i.pos < n {
    return nil, g.E("Nothing to drop")
  }

  i.pos += n
  return i, nil
}

func (i *VecIter) Dup(g *G) (Val, E) {
  out := *i
  return &out, nil
}

func (i *VecIter) Eq(g *G, rhs Val) bool {
  ri, ok := rhs.(*VecIter)
  return ok && ri.in.Eq(g, i.in) && ri.pos == i.pos
}

func (i *VecIter) Pop(g *G) (Val, Val, E) {
  if i.pos >= Int(len(i.in)) {
    return &g.NIL, i, nil
  }

  v := i.in[i.pos]
  i.pos++
  return v, i, nil
}

func (i *VecIter) Splat(g *G, out Vec) (Vec, E) {
  return append(out, i.in[i.pos:]), nil
}
