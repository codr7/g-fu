package gfu

import (
  "fmt"
)

type Error interface {
  String() string
}

type BasicError struct {
  pos Pos
  msg string
}

func (g *G) NewError(pos *Pos, msg string, args...interface{}) *BasicError {
  msg = fmt.Sprintf(msg, args...)

  if g.Debug {
    panic(msg)
  }
  
  return new(BasicError).Init(pos, msg)
}

func (e *BasicError) Init(pos *Pos, msg string) *BasicError {
  e.pos = *pos
  e.msg = msg
  return e
}

func (e *BasicError) String() string {
  p := &e.pos
  
  return fmt.Sprintf(
    "Error in '%s' on row %v, col %v:\n%v",
    p.src, p.row, p.col, e.msg)
}
