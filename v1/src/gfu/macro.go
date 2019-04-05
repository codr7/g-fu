package gfu

import (
)

type MacroImp func(*G, Pos, []Form, *Env) (Form, E)

type Macro struct {
  env *Env
  arg_list ArgList
  body []Form
  imp MacroImp
}

func NewMacro(g *G, env *Env, args []*Sym) *Macro {
  return new(Macro).Init(g, env, args)
}

func (m *Macro) Init(g *G, env *Env, args []*Sym) *Macro {
  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Macro) CallBody(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  if e := m.arg_list.CheckVals(g, pos, args); e != nil {
    return g.NIL, e
  }
  
  var be Env
  m.env.Clone(&be)
  m.arg_list.PutEnv(g, &be, args)
  return Forms(m.body).Eval(g, &be)
}

func (m *Macro) CallImp(g *G, pos Pos, args []Form, env *Env) (Form, E) {
  if e := m.arg_list.CheckForms(g, pos, args); e != nil {
    return nil, e
  }

  var f Form
  var e E
  
  if f, e = m.imp(g, pos, args, env); e != nil {
    return nil, e
  }

  return f, nil
}


func (e *Env) AddMacro(g *G, id string, imp MacroImp, args...string) {
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.S(a)
  }
  
  m := NewMacro(g, e, as)
  m.imp = imp
  
  var v Val
  v.Init(g.Macro, m)
  e.Put(g.S(id), v)
}
