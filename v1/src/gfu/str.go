package gfu

import (
  "bufio"
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

func (_ *StrType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteRune('"')

  for _, c := range val.(Str) {
    switch c {
    case '"':
      out.WriteString("\\\"")
    case '\x1b':
      out.WriteString("\\e")
    case '\n':
      out.WriteString("\\n")
    default:
      out.WriteRune(c)
    }
  }
  
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

func (_ *StrType) Print(g *G, val Val, out *bufio.Writer) E {
  out.WriteString(string(val.(Str)))
  return nil
}

func (g *G) String(val Val) (string, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)
  
  if e := g.Dump(val, w); e != nil {
    return "", e
  }

  w.Flush()
  return out.String(), nil
}

func (g *G) EString(val Val) string {
  s, e := g.String(val)

  if e != nil {
    s, _ = g.String(e)
  }

  return s
}
