package gfu

import (
  "bufio"
  "fmt"
)

type PrimImp func(*G, *Task, *Env, Vec, *Env) (Val, E)

type Prim struct {
  id       *Sym
  arg_list ArgList
  imp      PrimImp
}

type PrimType struct {
  BasicType
}

func NewPrim(g *G, id *Sym, imp PrimImp, args []Arg) *Prim {
  p := new(Prim)
  p.id = id
  p.arg_list.Init(g, args)
  p.imp = imp
  return p
}

func (_ *Prim) Type(g *G) Type {
  return &g.PrimType
}

func (_ *PrimType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (Val, E) {
  p := val.(*Prim)

  if e := p.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  return p.imp(g, task, env, p.arg_list.Fill(g, args), args_env)
}

func (_ *PrimType) Dump(g *G, val Val, out *bufio.Writer) E {
  p := val.(*Prim)
  fmt.Fprintf(out, "(prim %v", p.id)
  nargs := len(p.arg_list.items)

  if nargs > 0 {
    out.WriteString(" (")
  }

  for i, a := range p.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    if e := a.Dump(g, out); e != nil {
      return e
    }
  }

  if nargs > 0 {
    out.WriteRune(')')
  }

  out.WriteRune(')')
  return nil
}

func (env *Env) AddPrim(g *G, id string, imp PrimImp, args ...Arg) E {
  s := g.Sym(id)
  return env.Let(g, s, NewPrim(g, s, imp, args))
}

func ParsePrimArgs(g *G, args Val) Vec {
  if args == &g.NIL {
    return nil
  } else if v, ok := args.(Vec); ok {
    if len(v) == 0 {
      return nil
    }

    return v
  }

  return Vec{args}
}
