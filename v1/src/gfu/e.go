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

type EBasic struct {
  msg string
}

func (e *EBasic) Init(g *G, msg string) *EBasic {
  e.msg = msg
  return e
}

func (e *EBasic) Dump(g *G, out *bufio.Writer) E {
  out.WriteString(e.msg)
  return nil
}

func (e *EBasic) Msg(g *G) string {
  return e.msg
}

func (e *EBasic) Type(g *G) Type {
  return &g.EType
}

func (_ *EType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "(fail \"%v\")", val.(*EBasic).msg)
  return nil
}

func (_ *EType) Print(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("Error: ")
  out.WriteString(val.(*EBasic).msg)
  return nil
}

func (g *G) E(msg string, args ...interface{}) *EBasic {
  for i, a := range args {
    switch v := a.(type) {
    case Val:
      args[i] = g.EString(v)
    }
  }

  msg = fmt.Sprintf(msg, args...)
  e := new(EBasic).Init(g, msg)

  return e
}
