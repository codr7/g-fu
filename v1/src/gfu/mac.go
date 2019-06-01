package gfu

import (
  "bufio"
  "fmt"
  //"log"
)

type Mac struct {
  id       *Sym
  env      Env
  arg_list ArgList
  body     Vec
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
  
  m.arg_list.Init(g, args)
  return m, nil
}

func (m *Mac) InitEnv(g *G, env *Env) E {
  return g.Extenv(env, &m.env, m.body, false)
}

func (m *Mac) ExpandCall(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := make(Vec, len(args))
  var e E

  for i, a := range args {
    if avs[i], e = g.Quote(task, env, a, env); e != nil {
      return nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  var be Env
  m.env.Dup(&be)
  m.arg_list.LetVars(g, &be, args)
  return m.body.EvalExpr(g, task, &be, &be)
}

func (m *Mac) Type(g *G) Type {
  return &g.MacType
}

func (_ *MacType) ArgList(g *G, val Val) (*ArgList, E) {
  return &val.(*Mac).arg_list, nil
}

func (_ *MacType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (v Val, e E) {
  m := val.(*Mac)

  if v, e = m.ExpandCall(g, task, args_env, args); e != nil {
    return nil, e
  }

  if e = g.Extenv(&m.env, args_env, v, false); e != nil {
    return nil, e
  }

  if e = g.Extenv(&g.RootEnv, args_env, v, task != &g.MainTask); e != nil {
    return nil, e
  }

  return g.Eval(task, args_env, v, args_env)
}

func (_ *MacType) Dump(g *G, val Val, out *bufio.Writer) E {
  m := val.(*Mac)

  if id := m.id; id == nil {
    out.WriteString("(mac")
  } else {
    fmt.Fprintf(out, "(mac %v", m.id)
  }

  nargs := len(m.arg_list.items)

  if nargs > 0 {
    out.WriteString(" (")
  }

  for i, a := range m.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  if nargs > 0 {
    out.WriteRune(')')
  }

  out.WriteRune(')')
  return nil
}
