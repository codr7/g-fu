package gfu

import (
  //"log"
)

type Env struct {
  vars []Var
}

func (e *Env) Clone(dst *Env) {
  src := e.vars
  dst.vars = make([]Var, len(src))
  copy(dst.vars, src)
}

func (e *Env) Merge(keys []*Sym, vals []Val) {
  for i, k := range keys {
    v := vals[i]
    
    if j, found := e.Find(k); found == nil {
      e.Insert(j, k).Val = v
    } else {
      found.Val = v
    }
  }
}

func (e *Env) Find(key *Sym) (int, *Var) {
  vs := e.vars
  min, max := 0, len(vs)

  for min < max {
    i := (max-min)/2
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
  
  return vs[i].Init(key)
}

func (e *Env) Put(key *Sym, val_type Type, val interface{}) {
  i, found := e.Find(key)
  
  if found == nil {
    e.Insert(i, key).Val.Init(val_type, val)
  } else {
    found.Val.Init(val_type, val)
  }
}
