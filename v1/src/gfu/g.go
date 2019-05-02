package gfu

import (
  "fmt"
  "io/ioutil"
  //"log"
  "path/filepath"
  "strings"
  "sync"
)

type G struct {  
  Debug    bool
  MainTask Task
  RootEnv  Env

  MetaType,
  ChanType,
  FalseType, FunType,
  IntType, IterType,
  MacType,
  NilType,
  PrimType,
  QuoteType,
  SplatType, SpliceType, StrType, SymType,
  TaskType, TrueType,
  VecType Type

  NIL Nil
  T   True
  F   False

  syms,
  consts  sync.Map
  nsyms   uint64
  
  nil_sym *Sym
  load_path string
}

func NewG() (*G, E) {
  return new(G).Init()
}

func (g *G) Init() (*G, E) {
  g.nil_sym = g.Sym("_")
  g.MainTask.Init(g, &g.RootEnv, g.Sym("main-task"), NewChan(0), nil)
  return g, nil
}

func (g *G) NewEnv() *Env {
  var env Env
  g.RootEnv.Dup(g, &env)
  return &env
}

func (g *G) EvalString(task *Task, env *Env, pos Pos, s string) (Val, E) {
  in := strings.NewReader(s)
  var out Vec
  
  for {
    vs, e := g.Read(&pos, in, out, "")

    if e != nil {
      return nil, e
    }

    if vs == nil {
      break
    }

    out = vs
  }

  var e E
  if out, e = out.ExpandVec(g, task, env, -1); e != nil {
    return nil, e
  }
  
  return out.EvalExpr(g, &g.MainTask, env)
}

func (g *G) Load(task *Task, env *Env, path string) (Val, E) {
  path = filepath.Join(g.load_path, path)
  s, re := ioutil.ReadFile(path)

  if re != nil {
    return nil, g.E("Failed loading file: %v\n%v", path, re)
  }

  var pos Pos
  pos.Init(path)
  prev_path := g.load_path
  g.load_path = filepath.Dir(path)
  v, e := g.EvalString(task, env, pos, string(s))
  g.load_path = prev_path
  return v, e
}

func (e *Env) AddConst(g *G, id string, val Val) E {
  if _, dup := g.consts.LoadOrStore(g.Sym(id), val); dup {
    return g.E("Dup const: %v", id)
  }

  imp := func(g *G, task *Task, env *Env, args Vec) (Val, E) {
    v, e := args[0].Eval(g, task, env)

    if e != nil {
      return nil, e
    }
    
    if val.Eq(g, v) {
      return &g.T, nil
    }

    return &g.F, nil
  }

  return e.AddPrim(g, fmt.Sprintf("%v?", id), imp, A("val"))
}

func (g *G) FindConst(id *Sym) Val {
  v, ok := g.consts.Load(id)

  if !ok {
    return nil
  }

  return v.(Val)
}
