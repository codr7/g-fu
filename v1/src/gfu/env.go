package gfu

import (
//"log"
)

type Env struct {
  vars []Var
}

func (env *Env) Clone(g *G, dst *Env) (*Env, E)  {
  src := env.vars
  dst.vars = make([]Var, len(src))
  copy(dst.vars, src)
  var e E
  
  for i, _ := range dst.vars {
    v := &dst.vars[i]
    v.env = dst
    
    if v.Val, e = v.Val.Clone(g); e != nil {
      return nil, e
    }
  }

  return dst, nil
}

func (env *Env) Dup(g *G, dst *Env) (*Env, E) {
  src := env.vars
  dst.vars = make([]Var, len(src))
  copy(dst.vars, src)
  var e E
  
  for i, _ := range dst.vars {
    v := &dst.vars[i]
    
    if v.Val, e = v.Val.Dup(g); e != nil {
      return nil, e
    }
  }

  return dst, nil
}

func (e *Env) Find(key *Sym) (int, *Var) {
  vs := e.vars
  min, max := 0, len(vs)

  for min < max {
    i := (max + min) / 2
    v := &vs[i]

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

func (e *Env) Insert(i int, key *Sym) *Var {
  var v Var
  vs := append(e.vars, v)
  e.vars = vs

  if i < len(vs)-1 {
    copy(vs[i+1:], vs[i:])
  }

  return vs[i].Init(e, key)
}

func (e *Env) Let(key *Sym, val Val) {
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

  return found.Update(g, e, f)
}

type Var struct {
  env     *Env
  ext_var *Var
  key     *Sym
  Val     Val
}

func (v *Var) Init(env *Env, key *Sym) *Var {
  v.env = env
  v.key = key
  return v
}

func (v *Var) Update(g *G, env *Env, f func(Val) (Val, E)) (Val, E) {
  var e E

  if v.Val, e = f(v.Val); e != nil {
    return nil, e
  }

  if v.env != env {
    if v.ext_var == nil {
      _, v.ext_var = v.env.Find(v.key)
    }

    if v.ext_var == nil {
      return nil, g.E("Missing ext var: %v", v.key)
    }

    v.ext_var.Val = v.Val
  }

  return v.Val, nil
}
