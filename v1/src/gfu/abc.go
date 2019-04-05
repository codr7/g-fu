package gfu

import (
  //"log"
  "os"
  "strings"
  "time"
)

func do_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
  return Forms(args).Eval(g, env)
}

func fun_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
  asf := args[0]
  
  if _, ok := asf.(*ExprForm); !ok {
    return g.NIL, g.E(asf.Pos(), "Invalid fun args: %v", asf)
  }

  as, e := ArgsForm(asf.(*ExprForm).body).Parse(g)

  if e != nil {
    return g.NIL, e
  }

  f := NewFun(g, env, as)
  f.body = args[1:]
  
  var v Val
  v.Init(g.Fun, f)
  return v, nil
}

func let_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
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

func if_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
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

func and_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
  var e E
  v := g.NIL
  
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

func or_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
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

func for_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := Forms(args[1:])
  v := g.NIL
  
  for i := 0; i < n; i++ {
    if v, e = b.Eval(g, env); e != nil {
      return g.NIL, e
    }
  }
  
  return v, nil
}

func test_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
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

func bench_imp(g *G, pos Pos, args []Form, env *Env) (Val, E) {
  nv, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  n := nv.AsInt()
  b := Forms(args[1:])

  for i := 0; i < n; i++ {
    b.Eval(g, env)
  }

  t := time.Now()
  
  for i := 0; i < n; i++ {
    if _, e = b.Eval(g, env); e != nil {
      return g.NIL, e
    }
  }

  var v Val
  v.Init(g.Int, time.Now().Sub(t).Nanoseconds() / 1000000) 
  return v, nil
}

func dump_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  var out strings.Builder
  
  for _, v := range args {
    v.Dump(&out)
    out.WriteRune('\n')
  }

  os.Stderr.WriteString(out.String())
  return g.NIL, nil
}

func eval_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  f, e := args[0].Unquote(g, pos)

  if e != nil {
    return g.NIL, e
  }

  return f.Eval(g, env)
}

func recall_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  if g.recall_args != nil {
    return g.NIL, g.E(pos, "Recall already in progress")
  }
  
  g.recall_args = args
  return g.NIL, nil
}

func not_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  v := args[0]
  v.Init(g.Bool, !v.AsBool(g))
  return v, nil
}

func eq_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Eq(g, v) {
      v.Init(g.Bool, false)
      return v, nil
    }
  }
  
  v.Init(g.Bool, true)
  return v, nil
}

func is_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Is(g, v) {
      v.Init(g.Bool, false)
      return v, nil
    }
  }
  
  v.Init(g.Bool, true)
  return v, nil
}

func int_lt_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  v := args[0]
  a0 := v.AsInt()
  
  for _, a := range args[1:] {
    if a.AsInt() <= a0 {
      v.Init(g.Bool, false)
      return v, nil
    }
  }
  
  v.Init(g.Bool, true)
  return v, nil
}

func int_add_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  if len(args) == 1 {
    v := args[0]
    v.imp = Abs(v.AsInt())
    return v, nil
  }
  
  var v int

  for _, iv := range args {
    v += iv.AsInt()
  }
  
  var out Val
  out.Init(g.Int, v)
  return out, nil
}

func int_sub_imp(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  var out Val
  v := args[0].AsInt()

  if len(args) == 1 {
    out.Init(g.Int, -v)
  } else {    
    for _, iv := range args[1:] {
      v -= iv.AsInt()
    }
    
    out.Init(g.Int, v)
  }
  
  return out, nil
}

func (e *Env) InitAbc(g *G) {
  g.Bool = e.AddType(g, "Bool", new(BoolType))
  g.Fun = e.AddType(g, "Fun", new(FunType))
  g.Int = e.AddType(g, "Int", new(IntType))
  g.Macro = e.AddType(g, "Macro", new(MacroType))
  g.Meta = e.AddType(g, "Meta", new(MetaType))
  g.Nil = e.AddType(g, "Nil", new(NilType))
  g.Prim = e.AddType(g, "Prim", new(PrimType))
  g.Splat = e.AddType(g, "Splat", new(SplatType))
  g.Sym = e.AddType(g, "Sym", new(SymType))
  g.Vec = e.AddType(g, "Vec", new(VecType))
  
  e.AddVal(g, "_", g.Nil, nil, &g.NIL)
  e.AddVal(g, "T", g.Bool, true, &g.T)
  e.AddVal(g, "F", g.Bool, false, &g.F)
  
  e.AddPrim(g, "do", do_imp, "body..")
  e.AddPrim(g, "fun", fun_imp, "args", "body..")
  e.AddPrim(g, "let", let_imp, "args..")
  e.AddPrim(g, "if", if_imp, "cond", "t", "f?")
  e.AddPrim(g, "or", or_imp, "conds..")
  e.AddPrim(g, "and", and_imp, "conds..")
  e.AddPrim(g, "for", for_imp, "nreps", "body..")
  e.AddPrim(g, "test", test_imp, "cases..")
  e.AddPrim(g, "bench", bench_imp, "nreps", "body..")

  e.AddFun(g, "dump", dump_imp, "vals..")
  e.AddFun(g, "eval", eval_imp, "form")
  e.AddFun(g, "recall", recall_imp, "args..")

  e.AddFun(g, "not", not_imp, "x")
  
  e.AddFun(g, "=", eq_imp, "vals..")
  e.AddFun(g, "==", is_imp, "vals..")
  
  e.AddFun(g, "<", int_lt_imp, "vals..")
  e.AddFun(g, "+", int_add_imp, "vals..")
  e.AddFun(g, "-", int_sub_imp, "vals..")
}
