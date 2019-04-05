package gfu

import (
  "strings"
)

type MacroImp func(*G, Pos, []Form, *Env) (Form, E)

type Macro struct {
  env *Env

  min_args, max_args int
  args []*Sym

  body []Form
  imp MacroImp
}

func NewMacro(env *Env, args []*Sym) *Macro {
  return new(Macro).Init(env, args)
}

func (m *Macro) Init(env *Env, args []*Sym) *Macro {
  m.env = env
  m.args = args  
  nargs := len(args)

  if nargs > 0 {
    m.min_args, m.max_args = nargs, nargs
    a := args[nargs-1]
    
    if strings.HasSuffix(a.name, "..") {
      m.min_args--
      m.max_args = -1
    }
  }
  
  return m
}

func (e *Env) AddMacro(g *G, id string, imp MacroImp, args...string) {
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.S(a)
  }
  
  m := NewMacro(e, as)
  m.imp = imp
  
  var v Val
  v.Init(g.Macro, m)
  e.Put(g.S(id), v)
}
