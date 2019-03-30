package gfu

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
  avs, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  var out Val
  out.Init(g.Int, avs[0].AsInt() + avs[1].AsInt())
  return out, nil
}

func (e *Env) InitAbc(g *G) {
  g.Bool = new(BoolType).Init(g.Sym("Bool"))
  g.Fun = new(FunType).Init(g.Sym("Fun"))
  g.Int = new(IntType).Init(g.Sym("Int"))
  g.Nil = new(NilType).Init(g.Sym("Nil"))
  g.Prim = new(PrimType).Init(g.Sym("Prim"))
  
  e.AddVal(g, g.Sym("_"), g.Nil, nil, &g.NIL)
  e.AddVal(g, g.Sym("T"), g.Bool, true, &g.T)
  e.AddVal(g, g.Sym("F"), g.Bool, false, &g.F)
  
  e.AddPrim(g, g.Sym("bool"), 1, bool_imp)
  e.AddPrim(g, g.Sym("+"), 2, int_add_imp)
}
