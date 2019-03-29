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

func (t *FunType) Call(g *G, val Val, args []Val, env *Env) (Val, Error) {
  f := val.AsFun()
  
  if len(args) != len(f.args) {
    return g.NIL, g.NewError(&g.Pos, "Arg mismatch: %v", args)
  }
  
  out := g.NIL
  var be Env
  f.env.Clone(&be)
  be.Merge(f.args, args)
  
  for _, bf := range f.body {
    var e Error
    
    if out, e = bf.Eval(g, &be); e != nil {
      return g.NIL, g.NewError(&g.Pos, "Call failed: %v", e)
    }
  }

  return out, nil
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
