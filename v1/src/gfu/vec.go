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
  return v, nil
}

func (v Vec) Clone() Val {
  out := make(Vec, len(v))

  for i, it := range v {
    out[i] = it.Clone()
  }
  
  return out
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

func (v Vec) Eval(g *G, task *Task, env *Env) (Val, E) {
  if len(v) == 0 {
    return &g.NIL, nil
  }

  f := v[0]

  if s, ok := f.(*Sym); ok && s == g.nil_sym {
    return &g.NIL, nil
  }

  fv, e := f.Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  result, e := fv.Call(g, task, env, v[1:])

  if e != nil {
    return nil, g.E("Call failed: %v", e)
  }

  return result, nil
}

func (v Vec) EvalExpr(g *G, task *Task, env *Env) (Val, E) {
  var out Val = &g.NIL

  for _, it := range v {
    var e E

    if out, e = it.Eval(g, task, env); e != nil {
      return nil, e
    }

    if task.recall {
      break
    }
  }

  return out, nil
}

func (v Vec) EvalVec(g *G, task *Task, env *Env) (Vec, E) {
  var out Vec

  for _, it := range v {
    it, e := it.Eval(g, task, env)

    if e != nil {
      return nil, g.E("Arg eval failed: %v", e)
    }

    if task.recall {
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

func (v Vec) String() string {
  return DumpString(v)
}

func (v Vec) Type(g *G) *Type {
  return &g.VecType
}
