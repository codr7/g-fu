package gfu

import (
  "fmt"
)

type E interface {
  String() string
}

type BasicE struct {
  pos Pos
  msg string
}

func (e *BasicE) Init(pos Pos, msg string) *BasicE {
  e.pos = pos
  e.msg = msg
  return e
}

func (e *BasicE) String() string {
  p := &e.pos
  
  return fmt.Sprintf(
    "Error in '%s' on row %v, col %v:\n%v",
    p.src, p.Row, p.Col, e.msg)
}

func (g *G) E(pos Pos, msg string, args...interface{}) *BasicE {
  msg = fmt.Sprintf(msg, args...)  
  e := new(BasicE).Init(pos, msg)

  if g.Debug {
    panic(e.String())
  }

  return e
}
