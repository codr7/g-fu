package gfu

import (
  "strings"
)

type FunImp func(*G, Pos, []Val, *Env) (Val, E)

type Fun struct {
  min_args, max_args int
  args []*Sym
  body []Form
  env *Env
  imp FunImp
}

func NewFun(env *Env, args []*Sym) *Fun {
  return new(Fun).Init(env, args)
}

func (f *Fun) Init(env *Env, args []*Sym) *Fun {
  f.args = args
  f.env = env
  
  nargs := len(args)

  if nargs > 0 {
    f.min_args, f.max_args = nargs, nargs
    a := args[nargs-1]
    
    if strings.HasSuffix(a.name, "..") {
      f.min_args--
      f.max_args = -1
    }
  }
  
  return f
}

func (e *Env) AddFun(g *G, id string, imp FunImp, args...string) {
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.S(a)
  }
  
  f := NewFun(e, as)
  f.imp = imp
  
  var v Val
  v.Init(g.Fun, f)
  e.Put(g.S(id), v)
}
