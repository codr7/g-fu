package gfu

import (
  //"log"
  "strings"
)

type FunType struct {
  BasicType
}

func (t *FunType) Init(id *Sym) *FunType {
  t.BasicType.Init(id)
  return t
}

func (t *FunType) Call(g *G, val Val, args []Val, env *Env, pos Pos) (Val, Error) {
  f := val.AsFun()
  
  if len(args) != len(f.args) {
    return g.NIL, g.NewError(pos, "Arg mismatch: %v", args)
  }
  
  var be Env
  f.env.Clone(&be)
  be.Merge(f.args, args)
  return Forms(f.body).Eval(g, &be)
}

func (t *FunType) Dump(val Val, out *strings.Builder) {
  f := val.AsFun()
  out.WriteString("(fun (")

  for i, a := range f.args {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.name)
  }

  out.WriteString(") ")

  for i, bf := range f.body {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(bf.String())   
  }
  
  out.WriteRune(')')
}

func (v Val) AsFun() *Fun {
  return v.imp.(*Fun)
}
