package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Mac struct {
  BasicVal
  
  id       *Sym
  env      *Env
  env_cache Env
  arg_list ArgList
  body     Vec
}

func NewMac(g *G, env *Env, id *Sym, args []Arg) *Mac {
  return new(Mac).Init(g, env, id, args)
}

func (m *Mac) Init(g *G, env *Env, id *Sym, args []Arg) *Mac {
  m.BasicVal.Init(&g.MacType, m)

  if id != nil {
    m.id = id
    env.Let(id, m)
  }

  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Mac) ExpandCall(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := make(Vec, len(args))
  var e E

  for i, a := range args {
    if avs[i], e = a.Quote(g, task, env); e != nil {
      return nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  var be Env

  if m.env_cache.vars == nil {
    if e = m.body.Extenv(g, m.env, &be, false); e != nil {
      return nil, e
    }

    be.Dup(g, &m.env_cache)
  } else {
    m.env_cache.Dup(g, &be)
  }
  
  m.arg_list.LetVars(g, &be, args)
  return m.body.EvalExpr(g, task, &be)
}

func (m *Mac) Call(g *G, task *Task, env *Env, args Vec) (v Val, e E) {  
  if v, e = m.ExpandCall(g, task, env, args); e != nil {
    return nil, e
  }

  if e = v.Extenv(g, m.env, env, false); e != nil {
    return nil, e
  }

  return v.Eval(g, task, env)
}

func (m *Mac) Dump(out *strings.Builder) {
  if id := m.id; id == nil {
    out.WriteString("(mac (")
  } else {
    fmt.Fprintf(out, "(mac %v (", m.id)
  }
  
  for i, a := range m.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  out.WriteString(")")

  for _, bv := range m.body {
    bv.Dump(out)
  }

  out.WriteRune(')')
}
