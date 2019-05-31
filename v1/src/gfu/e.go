package gfu

import (
  "bufio"
  "fmt"
  //"log"
)

type E interface {
  Val
}

type EType struct {
  BasicType
}

type BasicE struct {
  msg string
}

func (e *BasicE) Init(g *G, msg string) *BasicE {
  if g.Debug {
    panic(msg)
  }

  e.msg = msg
  return e
}

func (e *BasicE) Dump(g *G, out *bufio.Writer) E {
  out.WriteString(e.msg)
  return nil
}

func (e *BasicE) Type(g *G) Type {
  return &g.EType
}

func (_ *EType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "(fail \"%v\")", val.(*BasicE).msg)
  return nil
}

func (_ *EType) Print(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("Error: ")
  out.WriteString(val.(*BasicE).msg)
  return nil
}

func (g *G) E(msg string, args ...interface{}) *BasicE {
  for i, a := range args {
    switch v := a.(type) {
    case Val:
      args[i] = g.EString(v)
    }
  }

  msg = fmt.Sprintf(msg, args...)
  e := new(BasicE).Init(g, msg)

  if g.Debug {
    panic(msg)
  }

  return e
}
