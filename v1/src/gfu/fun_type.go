package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunType struct {
  BasicType
}

func (t *FunType) Call(g *G, val Val, args Vec, env *Env) (Val, E) {
  f := val.AsFun()
  avs, e := args.EvalVec(g, env)

  if e != nil {
    return g.NIL, g.E("Args eval failed: %v", e)
  }

  if e := f.arg_list.Check(g, avs); e != nil {
    return g.NIL, e
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
    return g.NIL, e
  }
  
  if g.recall {
    avs, g.recall_args = g.recall_args, nil
    g.recall = false
    goto recall
  }
  
  return v, e
}

func (t *FunType) Dump(val Val, out *strings.Builder) {
  f := val.AsFun()
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
    
    for i, bf := range f.body {
      if i > 0 {
        out.WriteRune(' ')
      }
      
      out.WriteString(bf.String())   
    }
  
    out.WriteRune(')')
  }
}

func (v Val) AsFun() *Fun {
  return v.imp.(*Fun)
}
