package gfu

import (
  //"log"
  "os"
  "strings"
  "time"
)

func do_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args.EvalExpr(g, task, env)
}

func fun_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := args[0]
  
  if _, ok := avs.(Vec); !ok {
    return nil, g.E("Invalid args: %v", avs)
  }

  as, e := ParseArgs(g, task, env, avs.(Vec))

  if e != nil {
    return nil, e
  }

  f := NewFun(g, env, as)
  f.body = args[1:]
  return f, nil
}

func macro_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  avs := args[0]
  
  if _, ok := avs.(Vec); !ok {
    return nil, g.E("Invalid args: %v", avs)
  }

  as, e := ParseArgs(g, task, env, avs.(Vec))

  if e != nil {
    return nil, e
  }

  m := NewMacro(g, env, as)
  m.body = args[1:]
  return m, nil
}

func let_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  bsf := args[0]
  var bs Vec
  _, is_scope := bsf.(Vec)
  
  if is_scope {
    bs = bsf.(Vec)
  } else {
    bs = args
  }

  var le *Env

  if is_scope {
    le = new(Env)
    env.Clone(le)
  } else {
    le = env
  }
    
  for i := 0; i < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if _, ok := kf.(*Sym); !ok {
      return nil, g.E("Invalid let key: %v", kf)
    }

    k := kf.(*Sym)
    v, e := vf.Eval(g, task, le)

    if e != nil {
      return nil, e
    }

    le.Let(k, v)
  }

  if !is_scope {
    return &g.NIL, nil
  }
  
  rv, e := args[1:].EvalExpr(g, task, le)
  
  if e != nil {
    return nil, e
  }

  return rv, nil
}

func if_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  c, e := args[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  if c.Bool(g) {
    return args[1].Eval(g, task, env)
  }

  if len(args) > 2 {
    return args[2].Eval(g, task, env)
  }

  return &g.NIL, nil
}

func and_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var e E
  var v Val = &g.NIL
  
  for _, in := range args {
    v, e = in.Eval(g, task, env)

    if e != nil {
      return nil, e
    }
    
    if !v.Bool(g) {
      return &g.F, nil
    }
  }

  return v, nil
}

func or_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  for _, in := range args {
    v, e := in.Eval(g, task, env)

    if e != nil {
      return nil, e
    }
    
    if v.Bool(g) {
      return v, nil
    }
  }

  return &g.F, nil
}

func inc_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  d, e := args[1].Eval(g, task, env)

  if e != nil {
    return nil, e
  }
  
  return env.Update(g, args[0].(*Sym), func(v Val) (Val, E) {
    return v.(Int)+d.(Int), nil
  })
}

func for_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  nv, e := args[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  n := nv.(Int)
  b := args[1:]
  var v Val = &g.NIL
  
  for i := Int(0); i < n; i++ {
    if v, e = b.EvalExpr(g, task, env); e != nil {
      return nil, e
    }
  }
  
  return v, nil
}

func test_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  for _, in := range args {
    v, e := in.Eval(g, task, env)

    if e != nil {
      return nil, e
    }

    if !v.Bool(g) {
      return nil, g.E("Test failed: %v", in)
    }
  }

  return &g.NIL, nil
}

func bench_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  nv, e := args[0].Eval(g, task, env)

  if e != nil {
    return nil, e
  }

  n := nv.(Int)
  b := args[1:]

  for i := Int(0); i < n; i++ {
    b.EvalExpr(g, task, env)
  }

  t := time.Now()
  
  for i := Int(0); i < n; i++ {
    if _, e = b.EvalExpr(g, task, env); e != nil {
      return nil, e
    }
  }

  return Int(time.Now().Sub(t).Nanoseconds() / 1000000), nil
}

func dump_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder
  
  for _, v := range args {
    v.Dump(&out)
    out.WriteRune('\n')
  }

  os.Stderr.WriteString(out.String())
  return &g.NIL, nil
}

func eval_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Eval(g, task, env)
}

func recall_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if task.recall {
    return nil, g.E("Recall already in progress")
  }

  task.recall = true
  task.recall_args = args
  return &g.NIL, nil
}

func g_sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.GSym(""), nil
}

func bool_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if b := args[0].Bool(g); b {
    return &g.T, nil
  }

  return &g.F, nil
}

func not_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if b := args[0].Bool(g); b {
    return &g.F, nil
  }

  return &g.T, nil
}

func eq_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Eq(g, v) {
      return &g.F, nil
    }
  }
  
  return &g.T, nil
}

func is_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]
  
  for _, iv := range args[1:] {
    if !iv.Is(g, v) {
      return &g.F, nil
    }
  }
  
  return &g.T, nil
}

func int_lt_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  lhs := args[0].(Int)
  
  for _, a := range args[1:] {
    rhs := a.(Int)
    
    if rhs <= lhs {
      return &g.F, nil
    }

    lhs = rhs
  }
  
  return &g.T, nil
}

func int_add_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if len(args) == 1 {
    return args[0].(Int).Abs(), nil
  }
  
  var v Int

  for _, iv := range args {
    v += iv.(Int)
  }
  
  return v, nil
}

func int_sub_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0].(Int)

  if len(args) == 1 {
    return -v, nil
  }

  for _, iv := range args[1:] {
    v -= iv.(Int)
  }
  
  return v, nil
}

func vec_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args, nil
}

func vec_len_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Len(), nil
}

func vec_push_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return env.Update(g, args[0].(*Sym), func(v Val) (Val, E) {
    return v.(Vec).Push(args[1:]...), nil
  })
}

func vec_peek_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Peek(g), nil
}

func vec_pop_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var it Val
  
  env.Update(g, args[0].(*Sym), func(v Val) (Val, E) {
    var rest Vec
    it, rest = v.(Vec).Pop(g)
    return rest, nil
  })

  return it, nil
}

func task_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var e E
  as := args[0]
  var inbox Chan
  
  if v, ok := as.(Vec); ok {
    var a0 Val

    if a0, e = v[0].Eval(g, task, env); e != nil {
      return nil, e
    }

    if v, ok := a0.(Int); ok {
      inbox = NewChan(v)
    } else {
      inbox = a0.(Chan)
    }
  } else if s, ok := as.(*Sym); ok && s == g.nil_sym {
    inbox = NewChan(0)
  } else {
    return nil, g.E("Invalid task args: %v", as)
  }

  t := NewTask(g, inbox, args[1:])
  t.Start(g, env)
  return t, nil
}

func task_this_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return task, nil
}

func task_post_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  t := args[0]

  for _, v := range args[1:] {
    t.(*Task).Inbox<- v
  }
  
  return &g.NIL, nil
}

func task_fetch_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := <-task.Inbox

  if v == nil {
    v = &g.NIL
  }
  
  return v, nil
}

func task_wait_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  nargs := len(args)
  
  if nargs == 1 {
    return args[0].(*Task).Wait(), nil
  }

  out := make(Vec, nargs)
  
  for i, a := range args {
    out[i] = a.(*Task).Wait()
  }

  return out, nil
}

func chan_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return NewChan(args[0].(Int)), nil
}

func (e *Env) InitAbc(g *G) {
  e.AddType(g, &g.MetaType, "Meta")
  e.AddType(g, &g.ChanType, "Chan")
  e.AddType(g, &g.FalseType, "False")
  e.AddType(g, &g.FunType, "Fun")
  e.AddType(g, &g.IntType, "Int")
  e.AddType(g, &g.MacroType, "Macro")
  e.AddType(g, &g.NilType, "Nil")
  e.AddType(g, &g.PrimType, "Prim")
  e.AddType(g, &g.QuoteType, "Quote")
  e.AddType(g, &g.SpliceType, "Splice")
  e.AddType(g, &g.SplatType, "Splat")
  e.AddType(g, &g.SymType, "Sym")
  e.AddType(g, &g.TaskType, "Task")
  e.AddType(g, &g.TrueType, "True")
  e.AddType(g, &g.VecType, "Vec")
  
  e.AddVal(g, "_", &g.NIL)
  e.AddVal(g, "T", &g.T)
  e.AddVal(g, "F", &g.F)
  
  e.AddPrim(g, "do", do_imp, SplatA("body"))
  e.AddPrim(g, "fun", fun_imp, A("args"), SplatA("body"))
  e.AddPrim(g, "macro", macro_imp, A("args"), SplatA("body"))
  e.AddPrim(g, "let", let_imp, SplatA("args"))
  e.AddPrim(g, "if", if_imp, A("cond"), A("t"), OptA("f", nil))
  e.AddPrim(g, "or", or_imp, SplatA("conds"))
  e.AddPrim(g, "and", and_imp, SplatA("conds"))
  e.AddPrim(g, "inc", inc_imp, A("var"), OptA("delta", Int(1)))
  e.AddPrim(g, "for", for_imp, A("nreps"), SplatA("body"))
  e.AddPrim(g, "test", test_imp, SplatA("cases"))
  e.AddPrim(g, "bench", bench_imp, A("nreps"), SplatA("body"))

  e.AddFun(g, "dump", dump_imp, SplatA("vals"))
  e.AddFun(g, "eval", eval_imp, A("form"))
  e.AddFun(g, "recall", recall_imp, SplatA("args"))
  e.AddFun(g, "g-sym", g_sym_imp, OptA("prefix", nil))

  e.AddFun(g, "bool", bool_imp, A("val"))
  e.AddFun(g, "not", not_imp, A("val"))
  
  e.AddFun(g, "=", eq_imp, SplatA("vals"))
  e.AddFun(g, "==", is_imp, SplatA("vals"))
  
  e.AddFun(g, "<", int_lt_imp, SplatA("vals"))
  e.AddFun(g, "+", int_add_imp, SplatA("vals"))
  e.AddFun(g, "-", int_sub_imp, SplatA("vals"))

  e.AddFun(g, "vec", vec_imp, SplatA("items"))
  e.AddFun(g, "len", vec_len_imp, A("vec"))
  e.AddPrim(g, "push", vec_push_imp, A("vec"), SplatA("vals"))
  e.AddFun(g, "peek", vec_peek_imp, A("vec"))
  e.AddPrim(g, "pop", vec_pop_imp, A("vec"))

  e.AddPrim(g, "task", task_imp, A("args"), SplatA("body"))
  e.AddFun(g, "this-task", task_this_imp)
  e.AddFun(g, "post", task_post_imp, A("task"), SplatA("vals"))
  e.AddFun(g, "fetch", task_fetch_imp)
  e.AddFun(g, "wait", task_wait_imp, SplatA("tasks"))
  e.AddFun(g, "chan", chan_imp, OptA("buf", Int(0)))
}
