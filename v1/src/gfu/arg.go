package gfu

import (
  "fmt"
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
      a.id = g.Sym(idn[:len(idn)-1])
      l.min--
    } else if strings.HasSuffix(id.name, "..") {
      a.arg_type = ARG_SPLAT
      idn := id.name
      a.id = g.Sym(idn[:len(idn)-2])
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

func (l *ArgList) PutEnv(g *G, env *Env, args Vec) {
  nargs := len(args)
  
  for i, a := range l.items {
    if a.arg_type == ARG_SPLAT {
      var v Vec

      if i < nargs {
        v = make(Vec, nargs-i)
        copy(v, args[i:])
      }
      
      env.Put(a.id, v)
      break
    }

    if i < nargs {
      env.Put(a.id, args[i])
    } else {
      env.Put(a.id, &g.NIL)
    }
  }
}

type Args Vec

func (vs Args) Parse(g *G) ([]*Sym, E) {
  var out []*Sym
  
  for _, v := range vs {
    var id *Sym
    
    if sv, ok := v.(*Sym); ok {
      id = sv
    } else if ov, ok := v.(Opt); ok {
      id = g.Sym(fmt.Sprintf("%v?", ov.val.(*Sym)))
    } else if sv, ok := v.(Splat); ok {
      id = g.Sym(fmt.Sprintf("%v..", sv.val.(*Sym)))
    } else {
      return nil, g.E("Invalid arg: %v", v)
    }
    
    out = append(out, id)
  }

  return out, nil
}
