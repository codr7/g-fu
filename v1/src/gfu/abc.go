package gfu

import (
  //"log"
  "os"
  "strings"
  "time"
)

func do_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  return Forms(args).Eval(g, env)
}

func fun_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  asf := args[0]
  
  if _, ok := asf.(*ExprForm); !ok {
    return g.NIL, g.E(asf.Pos(), "Invalid fun args: %v", asf)
  }
  
  var as []*Sym
  
  for _, af := range asf.(*ExprForm).body {
    if _, ok := af.(*IdForm); !ok {
      return g.NIL, g.E(af.Pos(), "Invalid fun arg: %v", af)
    }
    
    as = append(as, af.(*IdForm).id)
  }
  
  var fv Val
  fv.Init(g.Fun, NewFun(as, args[1:], env))
  return fv, nil
}

func let_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  bsf := args[0]

  if _, ok := bsf.(*ExprForm); !ok {
    return g.NIL, g.E(bsf.Pos(), "Invalid let bindings: %v", bsf)
  }

  bs := bsf.(*ExprForm).body
  var le Env
  env.Clone(&le)
  
  for i := 0; i < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if _, ok := kf.(*IdForm); !ok {
      return g.NIL, g.E(kf.Pos(), "Invalid let key: %v", kf)
    }

    k := kf.(*IdForm).id
    v, e := vf.Eval(g, &le)

    if e != nil {
      return g.NIL, e
    }

    le.Put(k, v)
  }

  if len(args) == 1 {
    return g.NIL, nil
  }
  
  rv, e := Forms(args[1:]).Eval(g, &le)
  
  if e != nil {
    return g.NIL, e
  }
  
  return rv, nil
}

func if_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  c, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  if c.AsBool(g) {
    return args[1].Eval(g, env)
  }

  if len(args) > 2 {
    return args[2].Eval(g, env)
  }

  return g.NIL, nil
}

func and_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  var e Error
  var v Val
  
  for _, in := range args {
    v, e = in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }
    
    if !v.AsBool(g) {
      return g.F, nil
    }
  }

  return v, nil
}

func or_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  for _, in := range args {
    v, e := in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }
    
    if v.AsBool(g) {
      return v, nil
    }
  }

  return g.F, nil
}

func not_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  v, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  v.Init(g.Bool, !v.AsBool(g))
  return v, nil
}

func for_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := Forms(args[1:])
  v := g.NIL
  
  for i := Int(0); i < n; i++ {
    if v, e = b.Eval(g, env); e != nil {
      return g.NIL, e
    }
  }
  
  return v, nil
}

func recall_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  if g.recall_args != nil {
    return g.NIL, g.E(pos, "Recall already in progress")
  }
  
  var e Error
  
  if g.recall_args, e = args.Eval(g, env); e != nil {
    return g.NIL, g.E(pos, "Recall failed: %v", e)
  }
  
  return g.NIL, nil
}

func dump_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  var out strings.Builder
  
  for _, in := range args {
    v, e := in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }

    v.Dump(&out)
    out.WriteRune('\n')
  }

  os.Stderr.WriteString(out.String())
  return g.NIL, nil
}

func test_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  for _, in := range args {
    v, e := in.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }

    if !v.AsBool(g) {
      return g.NIL, g.E(pos, "Test failed")
    }
  }

  return g.NIL, nil
}

func bench_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := Forms(args[1:])

  for i := Int(0); i < n; i++ {
    b.Eval(g, env)
  }

  t := time.Now()
  
  for i := Int(0); i < n; i++ {
    if _, e = b.Eval(g, env); e != nil {
      return g.NIL, e
    }
  }

  var v Val
  v.Init(g.Int, time.Now().Sub(t).Nanoseconds() / 1000000) 
  return v, nil
}

func as_bool_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  var out Val
  out.Init(g.Bool, in.AsBool(g))
  return out, nil
}

func eq_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args.Eval(g, env)
  
  if e != nil {
    return g.NIL, e
  }

  if e = g.prim.CheckArgs(g, pos, in); e != nil {
    return g.NIL, e
  }
    
  var out Val
  v := in[0]
  
  for _, iv := range in[1:] {
    if !iv.Eq(g, v) {
      out.Init(g.Bool, false)
      return out, nil
    }
  }
  
  out.Init(g.Bool, true)
  return out, nil
}

func is_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args.Eval(g, env)
  
  if e != nil {
    return g.NIL, e
  }

  if e = g.prim.CheckArgs(g, pos, in); e != nil {
    return g.NIL, e
  }
    
  var out Val
  v := in[0]
  
  for _, iv := range in[1:] {
    if !iv.Is(g, v) {
      out.Init(g.Bool, false)
      return out, nil
    }
  }
  
  out.Init(g.Bool, true)
  return out, nil
}

func int_lt_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  if e = g.prim.CheckArgs(g, pos, in); e != nil {
    return g.NIL, e
  }

  var out Val
  v := in[0].AsInt()
  
  for _, iv := range in[1:] {
    if iv.AsInt() <= v {
      out.Init(g.Bool, false)
      return out, nil
    }
  }
  
  out.Init(g.Bool, true)
  return out, nil
}

func int_add_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  if e = g.prim.CheckArgs(g, pos, in); e != nil {
    return g.NIL, e
  }

  var out Val
  var v Int
  
  for _, iv := range in {
    v += iv.AsInt()
  }
  
  out.Init(g.Int, v)
  return out, nil
}

func int_sub_imp(g *G, pos Pos, args ListForm, env *Env) (Val, Error) {
  in, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  if e = g.prim.CheckArgs(g, pos, in); e != nil {
    return g.NIL, e
  }

  var out Val

  if len(in) == 1 {
    out.Init(g.Int, -in[0].AsInt())
  } else {
    v := in[0].AsInt()
    
    for _, iv := range in[1:] {
      v -= iv.AsInt()
    }
    
    out.Init(g.Int, v)
  }
  
  return out, nil
}

func (e *Env) InitAbc(g *G) {
  g.Bool = e.AddType(g, new(BoolType).Init(g.S("Bool")))
  g.Fun = e.AddType(g, new(FunType).Init(g.S("Fun")))
  g.Int = e.AddType(g, new(IntType).Init(g.S("Int")))
  g.Meta = e.AddType(g, new(MetaType).Init(g.S("Meta")))
  g.Nil = e.AddType(g, new(NilType).Init(g.S("Nil")))
  g.Prim = e.AddType(g, new(PrimType).Init(g.S("Prim")))
  g.Splat = e.AddType(g, new(SplatType).Init(g.S("Splat")))
  g.Vec = e.AddType(g, new(VecType).Init(g.S("Vec")))
  
  e.AddVal(g, g.S("_"), g.Nil, nil, &g.NIL)
  e.AddVal(g, g.S("T"), g.Bool, true, &g.T)
  e.AddVal(g, g.S("F"), g.Bool, false, &g.F)
  
  e.AddPrim(g, g.S("do"), 0, -1, do_imp)
  e.AddPrim(g, g.S("fun"), 1, -1, fun_imp)
  e.AddPrim(g, g.S("let"), 1, -1, let_imp)
  e.AddPrim(g, g.S("if"), 2, 3, if_imp)
  e.AddPrim(g, g.S("or"), 1, -1, or_imp)
  e.AddPrim(g, g.S("and"), 1, -1, and_imp)
  e.AddPrim(g, g.S("not"), 1, 1, not_imp)
  e.AddPrim(g, g.S("for"), 1, -1, for_imp)
  e.AddPrim(g, g.S("recall"), 0, -1, recall_imp)
  e.AddPrim(g, g.S("dump"), 1, -1, dump_imp)
  e.AddPrim(g, g.S("test"), 1, -1, test_imp)
  e.AddPrim(g, g.S("bench"), 1, -1, bench_imp)

  e.AddPrim(g, g.S("as-bool"), 1, 1, as_bool_imp)
  
  e.AddPrim(g, g.S("="), 2, -1, eq_imp)
  e.AddPrim(g, g.S("=="), 2, -1, is_imp)

  e.AddPrim(g, g.S("<"), 2, -1, int_lt_imp)
  e.AddPrim(g, g.S("+"), 1, -1, int_add_imp)
  e.AddPrim(g, g.S("-"), 1, -1, int_sub_imp)
}
