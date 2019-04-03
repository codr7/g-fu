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

  prim *Prim
  recall_args []Val
  
  Bool, Fun, Int, Meta, Nil, Prim, Splat, Sym, Vec Type
  NIL, T, F Val
}

func NewG() (*G, Error) {
  return new(G).Init()
}

func (g *G) Init() (*G, Error) {
  g.syms = make(Syms)
  return g, nil
}

func (g *G) EvalString(pos Pos, s string, env *Env) (Val, Error) {
  in := strings.NewReader(s)
  var out Forms
  
  for {
    fs, e := g.Read(&pos, in, out, 0)
    
    if e != nil {
      return g.NIL, e
    }
    
    if fs == nil {
      break
    }

    out = fs
  }

  return out.Eval(g, env)  
}

func (g *G) Load(pos Pos, fname string, env *Env) (Val, Error) {
  s, e := ioutil.ReadFile(fname)
  
  if e != nil {
    return g.NIL, g.E(pos, "Error loading file: %v\n%v", fname, e)
  }

  var fpos Pos
  fpos.Init(fname)
  return g.EvalString(fpos, string(s), env)
}
