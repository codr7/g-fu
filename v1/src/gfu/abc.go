package gfu

func do_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  return Forms(args).Eval(g, env)
}

func fun_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  asf := args[0]
  
  if asf.FormType() != &FORM_EXPR {
    return g.NIL, g.NewError(asf.Pos(), "Invalid fun args: %v", asf)
  }
  
  var as []*Sym
  
  for _, af := range asf.(*ExprForm).body {
    if af.FormType() != &FORM_ID {
      return g.NIL, g.NewError(af.Pos(), "Invalid fun arg: %v", af)
    }
    
    as = append(as, af.(*IdForm).id)
  }
  
  var fv Val
  fv.Init(g.Fun, NewFun(as, args[1:], env))
  return fv, nil
}

func let(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  bsf := args[0]

  if bsf.FormType() != &FORM_EXPR {
    return g.NIL, g.NewError(bsf.Pos(), "Invalid let bindings: %v", bsf)
  }

  bs := bsf.(*ExprForm).body
  var le Env
  env.Clone(&le)

  for i := 0; i < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if kf.FormType() != &FORM_ID {
      return g.NIL, g.NewError(kf.Pos(), "Invalid let key: %v", kf)
    }

    k := kf.(*IdForm).id
    v, e := vf.Eval(g, env)

    if e != nil {
      return g.NIL, e
    }

    i, found := le.Find(k)

    if found == nil {
      le.Insert(i, k).Val =  v
    } else {
      found.Val = v
    }
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

func bool_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  in, e := args[0].Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  var out Val
  out.Init(g.Bool, in.AsBool(g))
  return out, nil
}

func int_add_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  in, e := args.Eval(g, env)

  if e != nil {
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

func int_sub_imp(g *G, args ListForm, env *Env, pos Pos) (Val, Error) {
  in, e := args.Eval(g, env)

  if e != nil {
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
  g.Bool = new(BoolType).Init(g.Sym("Bool"))
  g.Fun = new(FunType).Init(g.Sym("Fun"))
  g.Int = new(IntType).Init(g.Sym("Int"))
  g.Nil = new(NilType).Init(g.Sym("Nil"))
  g.Prim = new(PrimType).Init(g.Sym("Prim"))
  g.Splat = new(SplatType).Init(g.Sym("Splat"))
  g.Vec = new(VecType).Init(g.Sym("Vec"))
  
  e.AddVal(g, g.Sym("_"), g.Nil, nil, &g.NIL)
  e.AddVal(g, g.Sym("T"), g.Bool, true, &g.T)
  e.AddVal(g, g.Sym("F"), g.Bool, false, &g.F)
  
  e.AddPrim(g, g.Sym("do"), 0, -1, do_imp)
  e.AddPrim(g, g.Sym("fun"), 1, -1, fun_imp)
  e.AddPrim(g, g.Sym("let"), 1, -1, let)

  e.AddPrim(g, g.Sym("bool"), 1, 1, bool_imp)
  e.AddPrim(g, g.Sym("+"), 0, -1, int_add_imp)
  e.AddPrim(g, g.Sym("-"), 1, -1, int_sub_imp)
}
