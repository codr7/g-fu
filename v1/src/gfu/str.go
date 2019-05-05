package gfu

import (
  //"log"
  "strings"
)

type Str []rune

type StrType struct {
  BasicType
}

func (s Str) Len() Int {
  return Int(len(s))
}

func (_ Str) Type(g *G) Type {
  return &g.StrType
}

func (_ *StrType) Bool(g *G, val Val) (bool, E) {
  return val.(Str).Len() > 0, nil
}

func (_ *StrType) Drop(g *G, val Val, n Int) (Val, E) {
  s := val.(Str)
  sl := s.Len()

  if sl < n {
    return nil, g.E("Nothing to drop")
  }

  return s[:sl-n], nil
}

func (_ *StrType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('"')
  out.WriteString(string(val.(Str)))
  out.WriteRune('"')
  return nil
}

func (_ *StrType) Eq(g *G, lhs, rhs Val) (bool, E) {
  rs, ok := rhs.(Str)
  return ok && string(lhs.(Str)) == string(rs), nil
}

func (_ *StrType) Len(g *G, val Val) (Int, E) {
  return val.(Str).Len(), nil
}

func (_ *StrType) Print(g *G, val Val, out *strings.Builder) {
  out.WriteString(string(val.(Str)))
}

func (g *G) String(val Val) (string, E) {
  var out strings.Builder

  if e := g.Dump(val, &out); e != nil {
    return "", e
  }
  
  return out.String(), nil
}
