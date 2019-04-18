package gfu

import (
  "io/ioutil"
  //"log"
  "path/filepath"
  "strings"
  "sync"
)

type G struct {
  syms    sync.Map
  nsyms   uint64
  nil_sym *Sym
  load_path string
  
  Debug    bool
  MainTask Task
  RootEnv  Env

  MetaType,
  ChanType, FalseType, FunType, IntType, MacroType, NilType, PrimType,
  QuoteType, SplatType, SpliceType, StrType, SymType, TaskType, TrueType, VecType Type

  NIL Nil
  T   True
  F   False
}

func NewG() (*G, E) {
  return new(G).Init()
}

func (g *G) Init() (*G, E) {
  g.nil_sym = g.Sym("_")
  g.MainTask.Init(g, NewChan(0), true, nil)
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

  if e := out.ExpandVec(g, task, env, -1); e != nil {
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
