package gfu

import (
  //"log"
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

func (a Arg) String() string {
  var out strings.Builder
  out.WriteString(a.id.name)

  switch a.arg_type {
  case ARG_OPT:
    out.WriteRune('?')
  case ARG_SPLAT:
    out.WriteString("..")
  }
  
  return out.String()
}

type ArgList struct {
  items []Arg
  min, max int
}

func (l *ArgList) Init(g *G, args []Arg) *ArgList {
  nargs := len(args)
  
  if nargs == 0 {
    return l
  }
  
  l.items = args
  l.min, l.max = nargs, nargs
  
  for _, a := range args {
    if a.arg_type == ARG_OPT {
      l.min--
    } 
  }
  
  a := l.items[nargs-1]
  
  if a.arg_type == ARG_SPLAT {
    l.min--
    l.max = -1
  }

  return l
}

func (l *ArgList) Check(g *G, args Vec) E {
  nargs := len(args)

  if (l.min != -1 && nargs < l.min) || (l.max != -1 && nargs > l.max) {
    return g.E("Arg mismatch")
  }

  return nil
}

func (l *ArgList) LetEnv(g *G, env *Env, args Vec) {
  nargs := len(args)
  
  for i, a := range l.items {
    if a.arg_type == ARG_SPLAT {
      var v Vec

      if i < nargs {
        v = make(Vec, nargs-i)
        copy(v, args[i:])
      }
      
      env.Let(a.id, v)
      break
    }

    if i < nargs {
      env.Let(a.id, args[i])
    } else {
      env.Let(a.id, &g.NIL)
    }
  }
}

func ParseArgs(g *G, in Vec) ([]Arg, E) {
  var out []Arg
  
  for _, v := range in {
    var a Arg
    
    if id, ok := v.(*Sym); ok {
      idn := id.name

      if strings.HasSuffix(idn, "?") {
        a.arg_type = ARG_OPT
        a.id = g.Sym(idn[:len(idn)-1])
      } else if strings.HasSuffix(idn, "..") {
        a.arg_type = ARG_SPLAT
        a.id = g.Sym(idn[:len(idn)-2])
      } else {
        a.id = id
      }
    } else if ov, ok := v.(Opt); ok {
      a.arg_type = ARG_OPT
      a.id = ov.val.(*Sym)
    } else if sv, ok := v.(Splat); ok {
      a.arg_type = ARG_SPLAT
      a.id = sv.val.(*Sym)
    } else {
      return nil, g.E("Invalid arg: %v", v)
    }
    
    out = append(out, a)
  }

  return out, nil
}
