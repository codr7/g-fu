package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunImp func(*G, Vec, *Env) (Val, E)

type Fun struct {
  env *Env
  arg_list ArgList
  body Vec
  imp FunImp
}

func NewFun(g *G, env *Env, args []*Sym) *Fun {
  return new(Fun).Init(g, env, args)
}

func (f *Fun) Init(g *G, env *Env, args []*Sym) *Fun {
  f.env = env
  f.arg_list.Init(g, args)
  return f
}

func (_ *Fun) Bool(g *G) bool {
  return true
}

func (f *Fun) Call(g *G, args Vec, env *Env) (Val, E) {
  avs, e := args.EvalVec(g, env)

  if e != nil {
    return nil, g.E("Args eval failed: %v", e)
  }

  if e := f.arg_list.Check(g, avs); e != nil {
    return nil, e
  }

  if f.imp != nil {
    return f.imp(g, avs, env)
  }
  
  var be Env
  f.env.Clone(&be)
  var v Val
recall:
  f.arg_list.PutEnv(g, &be, avs)

  if v, e = f.body.EvalExpr(g, &be); e != nil {
    g.recall_args = nil
    g.recall = false
    return nil, e
  }
  
  if g.recall {
    avs, g.recall_args = g.recall_args, nil
    g.recall = false
    goto recall
  }
  
  return v, e
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
    fmt.Fprintf(out, ") %v)", f.imp)
  } else {
    out.WriteString(") ")
    
    for i, bv := range f.body {
      if i > 0 {
        out.WriteRune(' ')
      }

      bv.Dump(out)
    }
  
    out.WriteRune(')')
  }
}

func (f *Fun) Eq(g *G, rhs Val) bool {
  return f == rhs
}

func (f *Fun) Eval(g *G, env *Env) (Val, E) {
  return f, nil
}

func (f *Fun) Is(g *G, rhs Val) bool {
  return f == rhs
}

func (f *Fun) Quote(g *G, env *Env) (Val, E) {
  return f, nil
}

func (f *Fun) Splat(g *G, out Vec) Vec {
  return append(out, f)
}

func (f *Fun) Type(g *G) *Type {
  return &g.FunType
}

func (e *Env) AddFun(g *G, id string, imp FunImp, args...string) {
  as := make([]*Sym, len(args))

  for i, a := range args {
    as[i] = g.Sym(a)
  }
  
  f := NewFun(g, e, as)
  f.imp = imp
  e.Put(g.Sym(id), f)
}
