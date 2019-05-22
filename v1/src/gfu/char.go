package gfu

import (
  "bufio"
)

type Char rune

type CharType struct {
  BasicType
}

func (_ Char) Type(g *G) Type {
  return &g.CharType
}

func (_ *CharType) Dump(g *G, val Val, out *bufio.Writer) E {
  c := rune(val.(Char))
  
  switch c {
  case '\n':
    out.WriteString("\\n")
  default:
    out.WriteRune(c)
  }
  
  return nil
}

func (_ *CharType) Print(g *G, val Val, out *bufio.Writer) {
  out.WriteRune(rune(val.(Char)))
}
