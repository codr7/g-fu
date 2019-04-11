package gfu

import (
  //"log"
)

type Macro struct {
  env *Env
  arg_list ArgList
  body Vec
}

func NewMacro(g *G, env *Env, args []*Sym) *Macro {
  return new(Macro).Init(g, env, args)
}

func (m *Macro) Init(g *G, env *Env, args []*Sym) *Macro {
  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Macro) Call(g *G, pos Pos, args Vec, env *Env) (Val, E) {
  var e E
  
  if e = m.arg_list.Check(g, pos, args); e != nil {
    return g.NIL, e
  }
  
  var be Env
  m.env.Clone(&be)
  m.arg_list.PutEnv(g, pos, &be, args)
  return m.body.EvalExpr(g, pos, &be)
}
