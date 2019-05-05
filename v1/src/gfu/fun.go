package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunImp func(*G, *Task, *Env, Vec) (Val, E)

type Fun struct {
  id        *Sym
  env       *Env
  env_cache Env
  arg_list  ArgList
  body      Vec
  imp       FunImp
}

type FunType struct {
  BasicType
}

func NewFun(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
  return new(Fun).Init(g, env, id, args)
}

func (f *Fun) Init(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
  if id != nil {
    f.id = id

    if e := env.Let(g, id, f); e != nil {
      return nil, e
    }
  }

  f.env = env
  f.arg_list.Init(g, args)
  return f, nil
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

  if f.env_cache.vars == nil {
    if e = g.Extenv(f.env, &be, f.body, false); e != nil {
      return nil, e
    }

    be.Dup(&f.env_cache)
  } else {
    f.env_cache.Dup(&be)
  }

recall:
  f.arg_list.LetVars(g, &be, args)
  var v Val

  if v, e = f.body.EvalExpr(g, task, &be); e != nil {
    if r, ok := e.(Recall); ok {
      be.Clear()
      f.env_cache.Dup(&be)
      args = r.args
      goto recall
    }

    return nil, e
  }

  return v, e
}

func (f *Fun) Type(g *G) Type {
  return &g.FunType
}

func (_ *FunType) Call(g *G, task *Task, env *Env, val Val, args Vec) (Val, E) {
  f := val.(*Fun)
  args, e := args.EvalVec(g, task, env)

  if e != nil {
    return nil, e
  }

  return f.CallArgs(g, task, env, args)
}

func (_ *FunType) Dump(g *G, val Val, out *strings.Builder) E {
  f := val.(*Fun)
  
  if id := f.id; id == nil {
    out.WriteString("(fun (")
  } else {
    fmt.Fprintf(out, "(fun %v (", f.id)
  }

  for i, a := range f.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    if e := a.Dump(g, out); e != nil {
      return e
    }
  }

  if f.imp == nil {
    out.WriteString(")")

    for _, bv := range f.body {
      out.WriteRune(' ')

      if e := g.Dump(bv, out); e != nil {
        return e
      }
    }

    out.WriteRune(')')
  } else {
    out.WriteString(") n/a)")
  }

  return nil
}

func (env *Env) AddFun(g *G, id string, imp FunImp, args ...Arg) E {
  f, e := NewFun(g, env, g.Sym(id), args)

  if e != nil {
    return e
  }

  f.imp = imp
  return nil
}

type Recall struct {
  args Vec
}

func NewRecall(args Vec) (r Recall) {
  r.args = args
  return r
}

func (r Recall) Dump(g *G, out *strings.Builder) {
  out.WriteString("(recall")

  for _, a := range r.args {
    g.Dump(a, out)
  }

  out.WriteRune(')')
}

func (r Recall) String() string {
  return "Recall"
}
