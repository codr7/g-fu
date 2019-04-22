package gfu

import (
  //"log"
)

type Env struct {
  vars []*Var
}

func (env *Env) Dup(g *G, dst *Env) (*Env, E) {
  src := env.vars
  dst.vars = make([]*Var, len(src))
  copy(dst.vars, src)
  return dst, nil
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

func (e *Env) Get(g *G, key *Sym) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return nil, g.E("Unknown: %v", key)
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

func (e *Env) Let(key *Sym, val Val) {
  if i, found := e.Find(key); found == nil {
    e.Insert(i, key).Val = val
  } else if found.env != e {
    v := new(Var).Init(e, key)
    v.Val = val
    e.vars[i] = v
  } else {
    found.Val = val
  }
}

func (e *Env) Set(g *G, key *Sym, val Val) (Val, E) {
  _, found := e.Find(key)

  if found == nil {
    return nil, g.E("Unknown var: %v", key)
  }

  var prev Val
  prev, found.Val = found.Val, val
  return prev, nil
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

  return v.Val, nil
}
