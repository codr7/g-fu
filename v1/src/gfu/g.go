package gfu

import (
  "io/ioutil"
  //"log"
  "strings"
)

type Syms map[string]*Sym

type G struct {
  Debug bool
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  recall bool
  recall_args Vec
  
  MetaType,
  FalseType, FunType, IntType, MacroType, NilType, OptType, PrimType, QuoteType,
  SplatType, SpliceType, SymType, TrueType, VecType Type
  
  NIL Nil
  T True
  F False
}

func NewG() (*G, E) {
  return new(G).Init()
}

func (g *G) Init() (*G, E) {
  g.syms = make(Syms)
  return g, nil
}

func (g *G) EvalString(pos Pos, s string, env *Env) (Val, E) {
  in := strings.NewReader(s)
  var out Vec

  for {
    vs, e := g.Read(&pos, in, Vec(out), 0)
    
    if e != nil {
      return g.NIL, e
    }
    
    if vs == nil {
      break
    }

    out = vs
  }

  return out.EvalExpr(g, env)  
}

func (g *G) Load(fname string, env *Env) (Val, E) {
  s, e := ioutil.ReadFile(fname)
  
  if e != nil {
    return g.NIL, g.E("Failed loading file: %v\n%v", fname, e)
  }

  var pos Pos
  pos.Init(fname)
  return g.EvalString(pos, string(s), env)
}
