package gfu

import (
  "bufio"
  "fmt"
  //"log"
  "strings"
  "time"
)

func do_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return args.EvalExpr(g, task, env, args_env)
}

func fun_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)

  if ok {
    i++
  }

  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]), args_env)

  if e != nil {
    return nil, e
  }

  i++
  f := NewFun(g, env, id, as...)
  f.body = args[i:]

  if e = f.InitEnv(g, env); e != nil {
    return nil, e
  }

  if args_env != env {
    if e := g.Extenv(args_env, &f.env, f.body, false); e != nil {
      return nil, e
    }
  }

  if id != nil {
    if e = env.Let(g, id, f); e != nil {
      return nil, e
    }
  }

  return f, nil
}

func pun_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  f, e := fun_imp(g, task, env, args, args_env)
  
  if e != nil {
    return nil, e
  }

  f.(*Fun).pure = true
  return f, nil
}

func mac_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  i := 0
  id, ok := args[0].(*Sym)

  if ok {
    i++
  }

  as, e := ParseArgs(g, task, env, ParsePrimArgs(g, args[i]), args_env)

  if e != nil {
    return nil, e
  }

  i++
  m, e := NewMac(g, env, id, as)

  if e != nil {
    return nil, e
  }

  m.body = args[i:]

  if e = m.InitEnv(g, env); e != nil {
    return nil, e
  }

  if args_env != env {
    if e := g.Extenv(args_env, &m.env, m.body, false); e != nil {
      return nil, e
    }
  }

  return m, nil
}

func call_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  t, e := g.Eval(task, args_env, args[0], args_env)

  if e != nil {
    return nil, e
  }

  return g.Call(task, env, t, args[1:], env)
}

func let_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (v Val, e E) {
  if len(args) == 0 {
    return &g.NIL, nil
  }

  bsf := args[0]
  bs, is_scope := bsf.(Vec)

  if bsf == &g.NIL {
    bs = nil
    is_scope = true
  }

  var le *Env

  if is_scope {
    le = new(Env)

    if e = g.Extenv(env, le, args, false); e != nil {
      return nil, e
    }
  } else {
    bs = args
    le = env
  }

  if e = g.Extenv(args_env, le, args, false); e != nil {
    return nil, e
  }

  v = &g.NIL

  for i := 0; i+1 < len(bs); i += 2 {
    kf, vf := bs[i], bs[i+1]

    if _, ok := kf.(*Sym); !ok {
      return nil, g.E("Invalid let key: %v", kf)
    }

    k := kf.(*Sym)
    v, e = g.Eval(task, le, vf, args_env)

    if e != nil {
      return nil, e
    }

    if e = le.Let(g, k, v); e != nil {
      return nil, e
    }
  }

  if !is_scope {
    return v, nil
  }

  rv, e := args[1:].EvalExpr(g, task, le, args_env)

  if e != nil {
    return nil, e
  }

  return rv, nil
}

func val_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  if v, _, _, e = args[0].(*Sym).Lookup(g, task, env, env, true); e != nil {
    return nil, e
  }

  if v == nil { v = &g.NIL }
  return v, nil
}

func set_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (v Val, e E) {
  for i := 0; i+1 < len(args); i += 2 {
    var k Val
    k, v = args[i], args[i+1]
    
    if v, e = g.Eval(task, env, v, args_env); e != nil {
      return nil, e
    }

    if e = env.Set(g, task, k, v, args_env); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func use_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  prefix := args[0]
  var ss []string
  
  for _, k := range args[1:] {
    if prefix == &g.NIL {
      ss = append(ss, k.(*Sym).name)
    } else { 
      ss = append(ss, fmt.Sprintf("%v/%v", prefix.(*Sym), k))
    }
  }

  if e := env.Use(g, args_env, ss...); e != nil {
    return nil, e
  }
  
  return &g.NIL, nil
}

func env_this_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return env, nil
}

func if_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  c, e := g.Eval(task, env, args[0], args_env)

  if e != nil {
    return nil, e
  }

  v, e := g.Bool(c)

  if e != nil {
    return nil, e
  }

  if v {
    return g.Eval(task, args_env, args[1], args_env)
  }

  if len(args) > 2 {
    return g.Eval(task, args_env, args[2], args_env)
  }

  return &g.NIL, nil
}

func inc_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  d, e := g.Eval(task, args_env, args[1], args_env)

  if e != nil {
    return nil, e
  }

  p := args[0]

  switch p.(type) {
  case *Sym, Vec:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Add(v, d)
    }, args_env)
  }

  if p, e = g.Eval(task, args_env, p, args_env); e != nil {
    return nil, e
  }

  return g.Add(p, d)
}

func test_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  for _, in := range args {
    v, e := g.Eval(task, env, in, args_env)

    if e != nil {
      return nil, e
    }

    bv, e := g.Bool(v)

    if e != nil {
      return nil, e
    }

    if !bv {
      return nil, g.E("Test failed: %v", in)
    }
  }

  return &g.NIL, nil
}

func bench_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  as := ParsePrimArgs(g, args[0])

  if as == nil {
    return nil, g.E("Invalid bench args: %v", as)
  }

  a, e := g.Eval(task, args_env, as[0], args_env)

  if e != nil {
    return nil, e
  }

  n := a.(Int)
  b := args[1:]

  for i := Int(0); i < n; i++ {
    b.EvalExpr(g, task, env, args_env)
  }

  t := time.Now()

  for i := Int(0); i < n; i++ {
    if _, e = b.EvalExpr(g, task, env, args_env); e != nil {
      return nil, e
    }
  }

  return Int(time.Now().Sub(t).Nanoseconds() / 1000000), nil
}

func debug_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  g.Debug = true
  return &g.NIL, nil
}

func throw_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, args[0]
}

func fail_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E(string(args[0].(Str)))
}

func abort_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, Abort{}
}

func retry_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if task.try == nil {
    return nil, g.E("Retry outside of try")
  }
  
  return nil, Retry{}
}

func try_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (ev Val, ee E) {
  var rfs []*Fun
  
  if args[0] != &g.NIL {
    rs, ok := args[0].(Vec)
    
    if !ok {
      return nil, g.E("Invalid restarts: %v", args[0].Type(g))
    }
    
    for _, r := range rs {
      rv, ok := r.(Vec)
      
      if !ok {
        return nil, g.E("Invalid restart: %v", r)
      }
      
      fargs, e := ParseArgs(g, task, env, ParsePrimArgs(g, rv[1]), args_env)
      
      if e != nil {
        return nil, e
      }
      
      var id *Sym
      
      if id, ok = rv[0].(*Sym); !ok {
        return nil, g.E("Invalid restart id: %v", rv[0].Type(g))
      }
      
      f := NewFun(g, env, id, fargs...)
      f.body = rv[2:]
      rfs = append(rfs, f)
    }
  }

  return g.Try(task, env, args_env, func() (Val, E) {
    return args[1:].EvalExpr(g, task, env, args_env)
  }, rfs...)
}

func catch_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  prev := len(task.catch_q)
  defer func() { task.catch_q = task.catch_q[:prev] }()
  handlers, ok := args[0].(Vec)

  if !ok {
    return nil, g.E("Invalid handlers: %v", args[0].Type(g))
  }

  for _, h := range handlers {
    hv, ok := h.(Vec)

    if !ok {
      return nil, g.E("Invalid handler: %v", h.Type(g))
    }

    var t Type
    var a Arg
    
    if hv[0] == &g.NIL {
      a.Init(g.NewSym(""))
    } else {
      as, ok := hv[0].(Vec)
      
      if !ok {
        return nil, g.E("Invalid handler args: %v", hv[0].Type(g))
      }
      
      tv, e := g.Eval(task, env, as[0], args_env)

      if e != nil {
        return nil, e
      }
      
      if tv != &g.NIL {
        t, ok = tv.(Type)
      }

      var s *Sym
      
      if as[1] != &g.NIL {
        s = as[1].(*Sym)
      }
      
      a.Init(s)
    }
    
    f := NewFun(g, env, nil, a)
    f.body = hv[1:]
    var c Catch
    c.Init(t, f)
    task.catch_q = append(task.catch_q, c)
  }

  v, e := args[1:].EvalExpr(g, task, env, args_env)

  if e != nil {
    return nil, e
  }
      
  return v, e
}

func restart_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  try := task.try
  
  if try == nil {
    return nil, g.E("Restart outside of try")
  }

  id, ok := args[0].(*Sym)

  if !ok {
    return nil, g.E("Invalid restart id: %v", args[0].Type(g))
  }
  
  v, e := try.restarts.Get(g, task, id, env, false)

  if e != nil {
    return nil, e
  }

  r := v.(Restart)
  r.args = args[1:]
  return r, nil
}

func load_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Load(task, env, env, string(args[0].(Str)))
}

func dup_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Dup(args[0])
}

func clone_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Clone(args[0])
}

func type_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].Type(g), nil
}

func eval_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  var e E
  v := args[0]

  if v, e = g.Eval(task, args_env, v, args_env); e != nil {
    return nil, e
  }

  return g.Eval(task, env, v, args_env)
}

func expand_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  return g.Expand(task, env, args[1], args[0].(Int))
}

func recall_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return &g.NIL, NewRecall(args)
}

func new_sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.NewSym(string(args[0].(Str))), nil
}

func sym_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)
  
  for _, a := range args {
    g.Print(a, w)
  }

  w.Flush()
  return g.Sym(out.String()), nil
}

func str_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)

  for _, a := range args {
    g.Print(a, w)
  }

  w.Flush()
  return Str(out.String()), nil
}

func bool_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.BoolVal(args[0])
}

func float_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Float(args[0])
}

func int_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Int(args[0])
}

func eq_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    ok, e := g.Eq(iv, v)

    if e != nil {
      return nil, e
    }

    if !ok {
      return &g.F, nil
    }
  }

  return &g.T, nil
}

func is_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  for _, iv := range args[1:] {
    if !g.Is(iv, v) {
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

func int_gt_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  lhs := args[0].(Int)

  for _, a := range args[1:] {
    rhs := a.(Int)

    if rhs >= lhs {
      return &g.F, nil
    }

    lhs = rhs
  }

  return &g.T, nil
}

func add_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Abs(a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Add(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func sub_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Neg(a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Sub(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func mul_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    return g.Mul(a0, a0)
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Mul(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func div_imp(g *G, task *Task, env *Env, args Vec) (v Val, e E) {
  a0 := args[0]

  if len(args) == 1 {
    var x, y Float
    x.SetInt(1)
    y.SetInt(a0.(Int))
    x.Div(y)
    return x, nil
  }

  v = args[0]

  for _, a := range args[1:] {
    if v, e = g.Div(v, a); e != nil {
      return nil, e
    }
  }

  return v, nil
}

func iter_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Iter(args[0])
}

func push_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  place := args[0]
  vs, e := args[1:].EvalVec(g, task, args_env, args_env)

  if e != nil {
    return nil, e
  }

  switch p := place.(type) {
  case *Sym:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Push(v, vs...)
    }, args_env)
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }
  }

  return g.Push(place, vs...)
}

func pop_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  var it, rest Val
  place := args[0]
  var e E

  switch p := place.(type) {
  case *Sym:
    args_env.Update(g, task, p, func(v Val) (Val, E) {
      if it, rest, e = g.Pop(v); e != nil {
        return nil, e
      }

      return rest, nil
    }, args_env)

    return it, nil
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }

    if it, rest, e = g.Pop(place); e != nil {
      return nil, e
    }
  }

  return NewSplat(g, Vec{it, rest}), nil
}

func drop_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  place := args[0]
  var e E

  switch p := place.(type) {
  case *Sym:
    return args_env.Update(g, task, p, func(v Val) (Val, E) {
      return g.Drop(v, args[1].(Int))
    }, args_env)
  default:
    if place, e = g.Eval(task, args_env, place, args_env); e != nil {
      return nil, e
    }
  }

  return g.Drop(place, args[1].(Int))
}

func len_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Len(args[0])
}

func seq_join_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  i, e := g.Iter(args[0])

  if e != nil {
    return nil, e
  }

  var sep Val
  var out strings.Builder
  w := bufio.NewWriter(&out)
    
  for {    
    v, _, e := g.Pop(i)

    if e != nil {
      return nil, e
    }
    
    if v == &g.NIL {
      break
    }

    if sep == nil {
      sep = args[1]
    } else if sep != &g.NIL {
      g.Print(sep, w)
    }

    g.Print(v, w)
  }

  w.Flush()
  return Str(out.String()), nil
}

func index_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.Index(args[0], args[1:])
}

func set_index_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return g.SetIndex(args[1], args[2:], args[0].(Setter))
}

func vec_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args, nil
}

func vec_peek_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Peek(g), nil
}

func find_key_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  in, k := args[0].(Vec), args[1]

  for i := 0; i < len(in)-1; i += 2 {
    if in[i] == k {
      return in[i+1], nil
    }
  }

  return &g.NIL, nil
}

func pop_key_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  in, k := args[0], args[1]
  var e E

  if k, e = g.Eval(task, args_env, k, args_env); e != nil {
    return nil, e
  }

  if id, ok := in.(*Sym); ok {
    var v Val

    if _, e = args_env.Update(g, task, id, func(in Val) (Val, E) {
      var out Val

      if v, out, e = in.(Vec).PopKey(g, k.(*Sym)); e != nil {
        return nil, e
      }

      return out, nil
    }, args_env); e != nil {
      return nil, e
    }

    return v, nil
  }

  if in, e = g.Eval(task, args_env, in, args_env); e != nil {
    return nil, e
  }

  var v, out Val

  if v, out, e = in.(Vec).PopKey(g, k.(*Sym)); e != nil {
    return nil, e
  }

  return NewSplat(g, Vec{v, out}), nil
}

func head_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  switch v := v.(type) {
  case Vec:
    if len(v) == 0 {
      return &g.NIL, nil
    }

    return v[0], nil
  case *Nil:
    return &g.NIL, nil
  default:
    break
  }

  return nil, g.E("Invalid head target: %v", v)
}

func tail_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  v := args[0]

  switch v := v.(type) {
  case Vec:
    if len(v) < 2 {
      return Vec(nil), nil
    }

    return v[1:], nil
  default:
    break
  }

  return nil, g.E("Invalid tail target: %v", v)
}

func reverse_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return args[0].(Vec).Reverse(), nil
}

func new_bin_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return NewBin(int(args[0].(Int))), nil
}

func bin_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  b := NewBin(len(args))
  
  for i, v := range args {
    bv, ok := v.(Byte)

    if !ok {
      return nil, g.E("Expected Byte: %v", v.Type(g))
    }

    b[i] = byte(bv)
  }
  
  return b, nil
}

func task_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  id, ok := args[0].(*Sym)
  i := 0

  if ok {
    i++
  }

  as := ParsePrimArgs(g, args[i])
  var inbox Chan
  var e E

  if as == nil {
    inbox = NewChan(0)
  } else {
    var a Val

    if a, e = g.Eval(task, args_env, as[0], args_env); e != nil {
      return nil, e
    }

    if v, ok := a.(Int); ok {
      inbox = NewChan(v)
    } else if v, ok := a.(Chan); ok {
      inbox = v
    } else {
      return nil, g.E("Invalid task args: %v", as)
    }
  }

  i++
  t, e := NewTask(g, env, id, inbox, args[i:])

  if e != nil {
    return nil, e
  }

  if e = t.Start(g, env); e != nil {
    return nil, e
  }

  return t, nil
}

func this_task_imp(g *G, task *Task, env *Env, args Vec, args_env *Env) (Val, E) {
  return task, nil
}

func task_post_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  t := args[0].(*Task)
  t.Inbox.Push(g, args[1:]...)
  return t, nil
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
  e.AddPrim(g, "do", true, do_imp, ASplat("body"))
  e.AddPrim(g, "fun", false, fun_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "pun", false, pun_imp, AOpt("id", nil), A("args"), ASplat("body"))
  e.AddPrim(g, "mac", false, mac_imp, AOpt("id", nil), A("args"), ASplat("body"))

  e.AddType(g, &g.MetaType, "Meta")
  e.AddType(g, &g.NumType, "Num")

  e.AddType(g, &g.SeqType, "Seq")
  g.SeqType.Env().AddPun(g, "join", seq_join_imp, A("in"), A("sep"))

  e.AddType(g, &g.IterType, "Iter", &g.SeqType)

  e.AddType(g, &g.FunType, "Fun")
  e.AddType(g, &g.PunType, "Pun", &g.FunType)

  e.AddType(g, &g.AbortType, "Abort")
  e.AddType(g, &g.BinType, "Bin", &g.SeqType)
  e.AddType(g, &g.BinIterType, "BinIter", &g.SeqType)
  e.AddType(g, &g.ByteType, "Byte", &g.NumType)
  e.AddType(g, &g.ChanType, "Chan")
  e.AddType(g, &g.CharType, "Char")
  e.AddType(g, &g.EnvType, "Env")
  e.AddType(g, &g.FalseType, "False")
  e.AddType(g, &g.FloatType, "Float", &g.NumType)
  e.AddType(g, &g.IntType, "Int", &g.NumType)
  e.AddType(g, &g.IntIterType, "IntIter", &g.IterType)
  e.AddType(g, &g.MacType, "Mac")
  e.AddType(g, &g.NilType, "Nil")
  e.AddType(g, &g.PrimType, "Prim")
  e.AddType(g, &g.QuoteType, "Quote")
  e.AddType(g, &g.RecallType, "Recall")
  e.AddType(g, &g.RestartType, "Restart")
  e.AddType(g, &g.RetryType, "Retry")
  e.AddType(g, &g.SetterType, "Setter")
  e.AddType(g, &g.SpliceType, "Splice")
  e.AddType(g, &g.SplatType, "Splat")
  e.AddType(g, &g.StrType, "Str")
  e.AddType(g, &g.SymType, "Sym")
  e.AddType(g, &g.TaskType, "Task")
  e.AddType(g, &g.TrueType, "True")
  e.AddType(g, &g.VecType, "Vec", &g.SeqType)
  e.AddType(g, &g.VecIterType, "VecIter", &g.IterType)
  e.AddType(g, &g.WriterType, "Writer")

  e.AddType(g, &g.EType, "E")
  e.AddType(g, &g.EImpureType, "EImpure")
  e.AddType(g, &g.EReadType, "ERead")
  e.AddType(g, &g.EUnknownType, "EUnknown", &g.EType)
  
  e.AddConst(g, "_", &g.NIL)
  e.AddConst(g, "T", &g.T)
  e.AddConst(g, "F", &g.F)
  e.AddConst(g, "\\e", Char('\x1b'))
  e.AddConst(g, "\\n", Char('\n'))
  e.AddConst(g, "\\s", Char(' '))
  e.AddConst(g, "\\\"", Char('"'))

  e.AddPrim(g, "call", false, call_imp, A("target"), ASplat("args"))
  e.AddPrim(g, "let", false, let_imp, ASplat("args"))
  e.AddFun(g, "val", val_imp, A("key"))
  e.AddPrim(g, "set", false, set_imp, ASplat("args"))
  e.AddPrim(g, "use", true, use_imp, AOpt("prefix", nil), ASplat("ids"))
  g.EnvType.Env().AddPrim(g, "this", false, env_this_imp)
  e.AddPrim(g, "if", true, if_imp, A("cond"), A("t"), AOpt("f", nil))
  e.AddPrim(g, "inc", true, inc_imp, A("var"), AOpt("delta", Int(1)))
  e.AddPrim(g, "test", true, test_imp, ASplat("cases"))
  e.AddPrim(g, "bench", true, bench_imp, A("nreps"), ASplat("body"))

  e.AddFun(g, "debug", debug_imp)
  e.AddFun(g, "throw", throw_imp, A("val"))
  e.AddFun(g, "fail", fail_imp, A("reason"))
  e.AddPrim(g, "try", true, try_imp, A("restarts"), ASplat("body"))
  e.AddPrim(g, "catch", true, catch_imp, A("handlers"), ASplat("body"))
  g.abort_fun, _ = e.AddFun(g, "abort", abort_imp)
  g.retry_fun, _ = e.AddFun(g, "retry", retry_imp)
  e.AddFun(g, "restart", restart_imp, A("id"), ASplat("args"))
  e.AddFun(g, "load", load_imp, A("path"))

  e.AddPun(g, "dup", dup_imp, A("val"))
  e.AddPun(g, "clone", clone_imp, A("val"))
  e.AddPun(g, "type", type_imp, A("val"))
  e.AddPrim(g, "eval", true, eval_imp, A("expr"))
  e.AddFun(g, "expand", expand_imp, A("n"), A("expr"))
  e.AddFun(g, "recall", recall_imp, ASplat("args"))
  e.AddFun(g, "new-sym", new_sym_imp, AOpt("prefix", Str("")))
  e.AddPun(g, "sym", sym_imp, ASplat("args"))
  e.AddPun(g, "str", str_imp, ASplat("args"))

  e.AddPun(g, "bool", bool_imp, A("val"))
  e.AddPun(g, "float", float_imp, A("val"))
  e.AddPun(g, "int", int_imp, A("val"))

  e.AddPun(g, "=", eq_imp, ASplat("vals"))
  e.AddPun(g, "==", is_imp, ASplat("vals"))

  e.AddPun(g, "<", int_lt_imp, ASplat("vals"))
  e.AddPun(g, ">", int_gt_imp, ASplat("vals"))

  e.AddPun(g, "+", add_imp, A("x"), ASplat("ys"))
  e.AddPun(g, "/", div_imp, A("x"), ASplat("ys"))
  e.AddPun(g, "-", sub_imp, A("x"), ASplat("ys"))
  e.AddPun(g, "*", mul_imp, A("x"), ASplat("ys"))

  e.AddPun(g, "iter", iter_imp, A("val"))
  e.AddPrim(g, "push", true, push_imp, A("out"), ASplat("vals"))
  e.AddPrim(g, "pop", true, pop_imp, A("in"))
  e.AddPrim(g, "drop", true, drop_imp, A("in"), AOpt("n", Int(1)))
  e.AddPun(g, "len", len_imp, A("in"))
  e.AddPun(g, "#", index_imp, A("source"), ASplat("key"))
  e.AddFun(g, "set-#", set_index_imp, A("set"), A("dest"), ASplat("key"))
  
  e.AddPun(g, "vec", vec_imp, ASplat("vals"))
  e.AddPun(g, "peek", vec_peek_imp, A("vec"))
  e.AddPun(g, "find-key", find_key_imp, A("in"), A("key"))
  e.AddPrim(g, "pop-key", true, pop_key_imp, A("in"), A("key"))
  e.AddPun(g, "head", head_imp, A("vec"))
  e.AddPun(g, "tail", tail_imp, A("vec"))
  e.AddFun(g, "reverse", reverse_imp, A("vec"))

  e.AddPun(g, "new-bin", new_bin_imp, AOpt("len", Int(0)))
  e.AddPun(g, "bin", bin_imp, ASplat("vals"))
  
  e.AddPrim(g, "task", false, task_imp, A("args"), ASplat("body"))
  g.TaskType.Env().AddPrim(g, "this", false, this_task_imp)
  g.TaskType.Env().AddFun(g, "post", task_post_imp, A("task"), A("val0"), ASplat("vals"))
  e.AddFun(g, "fetch", task_fetch_imp)
  e.AddFun(g, "wait", task_wait_imp, ASplat("tasks"))
  e.AddPun(g, "chan", chan_imp, AOpt("buf", Int(0)))
}
