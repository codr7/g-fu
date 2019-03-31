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

func (t *FunType) Call(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error) {
  f := val.AsFun()
  avs, e := args.Eval(g, env)
  nargs := len(avs)
  
  if (f.min_args != -1 && nargs < f.min_args) ||
    (f.max_args != -1 && nargs > f.max_args) {
    return g.NIL, g.NewError(pos, "Arg mismatch")
  }

  if e != nil {
    return g.NIL, g.NewError(pos, "Args eval failed: %v", e)
  }

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
