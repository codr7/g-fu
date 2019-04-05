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

func (m *Macro) CallBody(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  nargs := len(args)
  
  if (m.min_args != -1 && nargs < m.min_args) ||
    (m.max_args != -1 && nargs > m.max_args) {
    return g.NIL, g.E(pos, "Arg mismatch")
  }
  
  var be Env
  m.env.Clone(&be)
    
  for i, a := range m.args {
    id := a.name
    
    if strings.HasSuffix(id, "..") {
      v := new(Vec)
      v.items = make([]Val, nargs-i)
      copy(v.items, args[i:])
        
      var vv Val
      vv.Init(g.Vec, v)
      be.Put(g.S(id[:len(id)-2]), vv)
      break
    }

    be.Put(a, args[i])
  }
    
  return Forms(m.body).Eval(g, &be)
}

func (m *Macro) CallImp(g *G, pos Pos, args []Form, env *Env) (Form, E) {
  nargs := len(args)
  
  if (m.min_args != -1 && nargs < m.min_args) ||
    (m.max_args != -1 && nargs > m.max_args) {
    return nil, g.E(pos, "Arg mismatch")
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
  
  m := NewMacro(e, as)
  m.imp = imp
  
  var v Val
  v.Init(g.Macro, m)
  e.Put(g.S(id), v)
}
