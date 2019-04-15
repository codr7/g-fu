package gfu

import (
  //"log"
  "strings"
)

type Macro struct {
  env *Env
  arg_list ArgList
  body Vec
}

func NewMacro(g *G, env *Env, args []Arg) *Macro {
  return new(Macro).Init(g, env, args)
}

func (m *Macro) Init(g *G, env *Env, args []Arg) *Macro {
  m.env = env
  m.arg_list.Init(g, args)
  return m
}

func (m *Macro) Bool(g *G) bool {
  return true
}

func (m *Macro) Call(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  avs := make(Vec, len(args))
  
  for i, a := range args {
    if avs[i], e = a.Quote(g, task, env); e != nil {
      return nil, e
    }
  }

  if e = m.arg_list.Check(g, args); e != nil {
    return nil, e
  }
  
  var be Env
  m.env.Clone(&be)
  m.arg_list.LetEnv(g, &be, args)

  if v, e = m.body.EvalExpr(g, task, &be); e != nil {
    return nil, e
  }
  
  return v.Eval(g, task, env)
}

func (m *Macro) Dump(out *strings.Builder) {
  out.WriteString("(macro (")

  for i, a := range m.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  out.WriteString(") ")
  
  for i, bv := range m.body {
    if i > 0 {
      out.WriteRune(' ')
    }

    bv.Dump(out)
  }
  
  out.WriteRune(')')
}

func (m *Macro) Eq(g *G, rhs Val) bool {
  return m == rhs
}

func (m *Macro) Eval(g *G, task *Task, env *Env) (Val, E) {
  return m, nil
}

func (m *Macro) Is(g *G, rhs Val) bool {
  return m == rhs
}

func (m *Macro) Quote(g *G, task *Task, env *Env) (Val, E) {
  return m, nil
}

func (m *Macro) Splat(g *G, out Vec) Vec {
  return append(out, m)
}

func (m *Macro) Type(g *G) *Type {
  return &g.MacroType
}

