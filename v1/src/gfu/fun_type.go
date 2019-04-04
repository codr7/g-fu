package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type FunType struct {
  BasicType
}

func (t *FunType) Call(g *G, pos Pos, val Val, args VecForm, env *Env) (Val, E) {
  f := val.AsFun()
  avs, e := args.Eval(g, env)
  nargs := len(avs)
  
  if (f.min_args != -1 && nargs < f.min_args) ||
    (f.max_args != -1 && nargs > f.max_args) {
    return g.NIL, g.E(pos, "Arg mismatch")
  }

  if e != nil {
    return g.NIL, g.E(pos, "Args eval failed: %v", e)
  }
recall:
  var v Val

  if f.imp == nil {
    var be Env
    f.env.Clone(&be)
    
    for i, a := range f.args {
      id := a.name
      
      if strings.HasSuffix(id, "..") {
        v := new(Vec)
        v.items = make([]Val, nargs-i)
        copy(v.items, avs[i:])
        var vv Val
        vv.Init(g.Vec, v)
        be.Put(g.S(id[:len(id)-2]), vv)
        break
      }
      
      be.Put(a, avs[i])
    }
    
    if v, e = Forms(f.body).Eval(g, &be); e != nil {
      g.recall_args = nil
      return g.NIL, e
    }

    if g.recall_args != nil {
      avs, g.recall_args = g.recall_args, nil
      goto recall
    }
  } else {
    if v, e = f.imp(g, pos, avs, env); e != nil {
      g.recall_args = nil
      return g.NIL, e
    }
  }
  
  return v, nil
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
