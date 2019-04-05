package gfu

import (
  "strings"
)

type ArgType int

const (
  ARG_PLAIN ArgType = 0
  ARG_OPT ArgType = 1
  ARG_SPLAT ArgType = 2
)

type Arg struct {
  arg_type ArgType
  id *Sym
}

func (a *Arg) Init(id *Sym) *Arg {
  a.id = id
  return a
}

type ArgList struct {
  items []Arg
  min, max int
}

func (l *ArgList) Init(g *G, args []*Sym) *ArgList {
  nargs := len(args)
  
  if nargs == 0 {
    return l
  }
  
  l.items = make([]Arg, nargs)
  l.min, l.max = nargs, nargs
  
  for i, id := range args {
    a := &l.items[i]
    a.Init(id)

    if strings.HasSuffix(id.name, "?") {
      a.arg_type = ARG_OPT
      idn := id.name
      a.id = g.S(idn[:len(idn)-1])
      l.min--
    } else if strings.HasSuffix(id.name, "..") {
      a.arg_type = ARG_SPLAT
      idn := id.name
      a.id = g.S(idn[:len(idn)-2])
    }
  }
  
  a := l.items[nargs-1]
  
  if a.arg_type == ARG_SPLAT {
    l.min--
    l.max = -1
  }

  return l
}

func (l *ArgList) CheckVals(g *G, pos Pos, args []Val) E {
  nargs := len(args)

  if (l.min != -1 && nargs < l.min) || (l.max != -1 && nargs > l.max) {
    return g.E(pos, "Arg mismatch")
  }

  return nil
}

func (l *ArgList) CheckForms(g *G, pos Pos, args []Form) E {
  nargs := len(args)

  if (l.min != -1 && nargs < l.min) || (l.max != -1 && nargs > l.max) {
    return g.E(pos, "Arg mismatch")
  }

  return nil
}

func (l *ArgList) PutEnv(g *G, env *Env, args []Val) {
  nargs := len(args)
  
  for i, a := range l.items {
    if a.arg_type == ARG_SPLAT {
      v := new(Vec)
      v.items = make([]Val, nargs-i)
      copy(v.items, args[i:])

      var vv Val
      vv.Init(g.Vec, v)
      env.Put(a.id, vv)
      break
    }
      
    env.Put(a.id, args[i])
  }

}
