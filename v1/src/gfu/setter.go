package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Setter func(Val) (Val, E)

type SetterType struct {
  BasicType
}

func (_ Setter) Type(g *G) Type {
  return &g.SetterType
}

func (_ *SetterType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (v Val, e E) {
  if v, e = g.Eval(task, env, args[0]); e != nil {
    return nil, e
  }

  return val.(Setter)(v)
}

func (_ *SetterType) Dump(g *G, val Val, out *strings.Builder) E {
  fmt.Fprintf(out, "set-%v", func(Val) (Val, E)(val.(Setter)))
  return nil
}
