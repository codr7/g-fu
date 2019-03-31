package gfu

import (
  "io/ioutil"
  "strings"
)

type Syms map[string]*Sym

type G struct {
  Debug bool
  RootEnv Env
  
  sym_tag Tag
  syms Syms

  Bool, Fun, Int, Meta, Nil, Prim, Splat, Vec Type
  NIL, T, F Val
}

func NewG() (*G, Error) {
  return new(G).Init()
}

func (g *G) Init() (*G, Error) {
  g.syms = make(Syms)
  return g, nil
}

func (g *G) Load(fname string, env *Env, pos Pos) (Val, Error) {
  s, e := ioutil.ReadFile(fname)
  
  if e != nil {
    return g.NIL, g.E(pos, "Error loading file: %v\n%v", fname, e)
  }

  in := strings.NewReader(string(s))
  var fs Forms
  pos.Init(fname)
  
  for {
    f, e := g.Read(in, &pos, 0)

    if e != nil {
      return g.NIL, e
    }

    if f == nil {
      break
    }
    
    fs = append(fs, f)
  }

  return fs.Eval(g, env)
}
