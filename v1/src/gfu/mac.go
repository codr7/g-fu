package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Mac struct {
  id        *Sym
  env       *Env
  env_cache Env
  arg_list  ArgList
  body      Vec
}

type MacType struct {
  BasicType
}

func NewMac(g *G, env *Env, id *Sym, args []Arg) (*Mac, E) {
  return new(Mac).Init(g, env, id, args)
}

func (m *Mac) Init(g *G, env *Env, id *Sym, args []Arg) (*Mac, E) {
  if id != nil {
    m.id = id

    if e := env.Let(g, id, m); e != nil {
      return nil, e
    }
  }

  m.env = env
  m.arg_list.Init(g, args)
  return m, nil
}

func (m *Mac) ExpandCall(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := make(Vec, len(args))
  var e E

  for i, a := range args {
    if avs[i], e = g.Quote(task, env, a); e != nil {
      return nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  var be Env

  if m.env_cache.vars == nil {
    if e = g.Extenv(m.env, &be, m.body, false); e != nil {
      return nil, e
    }

    be.Dup(&m.env_cache)
  } else {
    m.env_cache.Dup(&be)
  }

  m.arg_list.LetVars(g, &be, args)
  return m.body.EvalExpr(g, task, &be)
}

func (m *Mac) Type(g *G) Type {
  return &g.MacType
}

func (_ *MacType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (v Val, e E) {
  m := val.(*Mac)
  
  if v, e = m.ExpandCall(g, task, env, args); e != nil {
    return nil, e
  }

  if e = g.Extenv(m.env, env, v, false); e != nil {
    return nil, e
  }

  return g.Eval(task, env, v)
}

func (_ *MacType) Dump(g *G, val Val, out *strings.Builder) E {
  m := val.(*Mac)
  
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

  out.WriteString("))")
  return nil
}
