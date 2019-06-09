package gfu

import (
  "bufio"
  "fmt"
  //"log"
)

type Setter func(Val) (Val, E)

type SetterType struct {
  BasicType
  arg_list ArgList
}

func (_ Setter) Type(g *G) Type {
  return &g.SetterType
}

func (t *SetterType) Init(g *G, id *Sym, parents []Type) (e E) {
  if e = t.BasicType.Init(g, id, parents); e != nil {
    return e
  }

  t.arg_list.Init(g, []Arg{A("val")})
  return nil
}

func (t *SetterType) ArgList(g *G, _ Val) (*ArgList, E) {
  return &t.arg_list, nil
}

func (_ *SetterType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (v Val, e E) {
  if v, e = g.Eval(task, env, args[0], args_env); e != nil {
    return nil, e
  }

  return val.(Setter)(v)
}

func (_ *SetterType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "set-%v", (func(Val) (Val, E))(val.(Setter)))
  return nil
}
