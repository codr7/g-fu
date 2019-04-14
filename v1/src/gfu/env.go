package gfu

import (
  //"log"
)

type Env struct {
  vars []Var
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

func (v *Var) Update(env *Env, f func(Val) (Val, E)) (Val, E) {
  var e E
  
  if v.Val, e = f(v.Val); e != nil {
    return nil, e
  }

  if v.env != env {
    v.env.Put(v.key, v.Val)
  }

  return v.Val, nil
}

func (e *Env) Clone(dst *Env) {
  src := e.vars
  dst.vars = make([]Var, len(src))
  copy(dst.vars, src)
}

func (e *Env) Find(key *Sym) (int, *Var) {
  vs := e.vars
  min, max := 0, len(vs)

  for min < max {
    i := (max+min)/2
    v := &vs[i]
    
    switch key.tag.Cmp(v.key.tag) {
    case -1:
      max = i
    case 1:
      min = i+1
    default:
      return i, v
    }
  }
  
  return max, nil
}

func (e *Env) Insert(i int, key *Sym) *Var {
  var v Var
  vs := append(e.vars, v)
  e.vars = vs

  if i < len(vs)-1 {
    copy(vs[i+1:], vs[i:])
  }
  
  return vs[i].Init(e, key)
}

func (e *Env) Put(key *Sym, val Val) {
  i, found := e.Find(key)
  
  if found == nil {
    e.Insert(i, key).Val = val
  } else {
    found.env = e
    found.Val = val
  }
}

func (e *Env) Update(g *G, key *Sym, f func(Val) (Val, E)) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return nil, g.E("Unknown var: %v", key)
  }

  return found.Update(e, f)
}
