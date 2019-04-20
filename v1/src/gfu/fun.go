package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunImp func(*G, *Task, *Env, Vec) (Val, E)

type Fun struct {
  BasicVal
  
  env      *Env
  arg_list ArgList
  body     Vec
  imp      FunImp
}

func NewFun(g *G, env *Env, args []Arg) *Fun {
  return new(Fun).Init(g, env, args)
}

func (f *Fun) Init(g *G, env *Env, args []Arg) *Fun {
  f.BasicVal.Init(&g.FunType, f)
  f.env = env
  f.arg_list.Init(g, args)
  return f
}

func (f *Fun) CallArgs(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var e E
  
  if e = f.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  if f.imp != nil {
    return f.imp(g, task, env, f.arg_list.Fill(g, args))
  }

  var be Env

  if e = f.body.Extenv(g, f.env, &be, false); e != nil {
    return nil, e
  }

  var v Val
recall:
  f.arg_list.LetVars(g, &be, args)

  if v, e = f.body.EvalExpr(g, task, &be); e != nil {
    task.recall_args = nil
    task.recall = false
    return nil, e
  }

  if task.recall {
    args, task.recall_args = task.recall_args, nil
    task.recall = false
    goto recall
  }

  return v, e
}

func (f *Fun) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  args, e := args.EvalVec(g, task, env)

  if e != nil {
    return nil, g.E("Args eval failed: %v", e)
  }

  return f.CallArgs(g, task, env, args)
}

func (f *Fun) Dump(out *strings.Builder) {
  out.WriteString("(fun (")

  for i, a := range f.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  if f.imp == nil {
    out.WriteString(") ")

    for i, bv := range f.body {
      if i > 0 {
        out.WriteRune(' ')
      }

      bv.Dump(out)
    }

    out.WriteRune(')')
  } else {
    fmt.Fprintf(out, ") %v)", f.imp)
  }
}

func (env *Env) AddFun(g *G, id string, imp FunImp, args ...Arg) E {
  f := NewFun(g, env, args)
  f.imp = imp
  env.Let(g.Sym(id), f)
  return nil
}
