package gfu

import (
  "strings"
)

type ArgList struct {
  items []*Sym
  min, max int
}

func (l *ArgList) Init(args []*Sym) *ArgList {
  l.items = args  
  nargs := len(args)

  if nargs > 0 {
    l.min, l.max = nargs, nargs
    a := args[nargs-1]
    
    if strings.HasSuffix(a.name, "..") {
      l.min--
      l.max = -1
    }
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
