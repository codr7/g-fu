package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Env struct {
  vars []*Var
  resolve *Fun
}

type EnvType struct {
  BasicType
}

func (e *Env) Clear() {
  e.vars = nil
}

func (e *Env) Dump(g *G, out *strings.Builder) E {
  out.WriteRune('(')

  for i, v := range e.vars {
    if i > 0 {
      out.WriteRune(' ')
    }

    fmt.Fprintf(out, "%v:", v.key)

    if v.Val == e {
      out.WriteString("this-form")
    } else if e := g.Dump(v.Val, out); e != nil {
      return e
    }
  }

  out.WriteRune(')')
  return nil
}

func (e *Env) Dup(dst *Env) *Env {
  src := e.vars
  dst.vars = make([]*Var, len(src))
  copy(dst.vars, src)
  return dst
}

func (dst *Env) Extend(g *G, src *Env, clone bool, keys...*Sym) E {
  for _, k := range keys {
    if i, dv := dst.Find(k); dv == nil {
      if _, sv := src.Find(k); sv != nil {
        if clone {
          dv = dst.Insert(i, sv.key)
          var e E
          
          if dv.Val, e = g.Clone(sv.Val); e != nil {
            return e
          }
        } else {
          dst.InsertVar(i, sv)
        }
      }
    }
  }

  return nil
}

func (e *Env) Find(key *Sym) (int, *Var) {
  vs := e.vars
  min, max := 0, len(vs)

  for min < max {
    i := (max + min) / 2
    v := vs[i]

    switch key.tag.Cmp(v.key.tag) {
    case -1:
      max = i
    case 1:
      min = i + 1
    default:
      return i, v
    }
  }

  return max, nil
}

func (e *Env) Get(g *G, task *Task, key *Sym, silent bool) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return e.Resolve(g, task, key, silent)
  }

  return found.Val, nil
}

func (e *Env) Insert(i int, key *Sym) *Var {
  v := new(Var).Init(e, key)
  e.InsertVar(i, v)
  return v
}

func (e *Env) InsertVar(i int, v *Var) {
  vs := e.vars
  vs = append(vs, v)

  if i < len(vs)-1 {
    copy(vs[i+1:], vs[i:])
    vs[i] = v
  }

  e.vars = vs
}

func (e *Env) Len() Int {
  return Int(len(e.vars))
}

func (e *Env) Let(g *G, key *Sym, val Val) E {
  if key == g.resolve_sym {
    var ok bool
    
    if e.resolve, ok = val.(*Fun); !ok {
      return g.E("Expected Fun, was: %v", val.Type(g))
    }
    
    return nil
  }
  
  if i, found := e.Find(key); found == nil {
    e.Insert(i, key).Val = val
  } else if found.env != e {
    v := new(Var).Init(e, key)
    v.Val = val
    e.vars[i] = v
  } else {
    return g.E("Dup binding: %v %v", key, found.Val)
  }

  return nil
}

func (e *Env) Resolve(g *G, task *Task, key *Sym, silent bool) (Val, E) {
  if e.resolve == nil {
    if silent {
      return nil, nil
    }
    
    return nil, g.E("Unknown: %v", key)
  }

  return e.resolve.CallArgs(g, task, e, Vec{key})
}

func (e *Env) Set(g *G, key *Sym, val Val) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return nil, g.E("Unknown: %v", key)
  }

  var prev Val
  prev, found.Val = found.Val, val
  return prev, nil
}

func (e *Env) Type(g *G) Type {
  return &g.EnvType
}

func (e *Env) Update(g *G, key *Sym, f func(Val) (Val, E)) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return nil, g.E("Unknown: %v", key)
  }

  return found.Update(g, e, f)
}

func (_ *EnvType) Bool(g *G, val Val) (bool, E) {
  return len(val.(*Env).vars) > 0, nil
}

func (_ *EnvType) Clone(g *G, val Val) (Val, E) {
  env := val.(*Env)
  src := env.vars
  dst := new(Env)
  dst.vars = make([]*Var, len(src))
  var e E
  
  for i, v := range src {
    if dst.vars[i], e = v.Clone(g, env); e != nil {
      return nil, e
    }
  }

  return dst, nil
}

func (_ *EnvType) Dump(g *G, val Val, out *strings.Builder) E {
  return val.(*Env).Dump(g, out)
}

func (_ *EnvType) Dup(g *G, val Val) (Val, E) {
  return val.(*Env).Dup(new(Env)), nil
}

func (_ *EnvType) Len(g *G, val Val) (Int, E) {
  return val.(*Env).Len(), nil
}

type Var struct {
  env *Env
  key *Sym
  Val Val
}

func (v *Var) Init(env *Env, key *Sym) *Var {
  v.env = env
  v.key = key
  return v
}

func (v *Var) Clone(g *G, env *Env) (dst *Var, e E) {
  dst = new(Var).Init(env, v.key)

  if dst.Val, e = g.Clone(v.Val); e != nil {
    return nil, e
  }

  return dst, e
}

func (v *Var) Dump(g *G, out *strings.Builder) {
  fmt.Fprintf(out, "%v:", v.key)
  g.Dump(v.Val, out)
}

func (v *Var) Update(g *G, env *Env, f func(Val) (Val, E)) (Val, E) {
  var e E

  if v.Val, e = f(v.Val); e != nil {
    return nil, e
  }

  return v.Val, nil
}
