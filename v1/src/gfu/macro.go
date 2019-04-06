package gfu

import (
  //"log"
)

type Macro struct {
  env *Env
  arg_list ArgList
  body []Form
}

func NewMacro(g *G, env *Env, args []*Sym) *Macro {
  return new(Macro).Init(g, env, args)
}

func (m *Macro) Init(g *G, env *Env, args []*Sym) *Macro {
  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Macro) Call(g *G, pos Pos, args []Val, env *Env) (Form, E) {
  var e E
  
  if e = m.arg_list.CheckVals(g, pos, args); e != nil {
    return nil, e
  }
  
  var be Env
  m.env.Clone(&be)
  m.arg_list.PutEnv(g, &be, args)
  var v Val
  
  if v, e = Forms(m.body).Eval(g, &be); e != nil {
    return nil, e
  }

  return v.Unquote(g, pos)
}
