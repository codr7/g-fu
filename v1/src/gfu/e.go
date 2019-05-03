package gfu

import (
	"fmt"
	//"log"
)

type E interface {
	String() string
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

func (e *BasicE) String() string {
	return e.msg
}

func (g *G) E(msg string, args ...interface{}) *BasicE {
	msg = fmt.Sprintf(msg, args...)
	e := new(BasicE).Init(g, fmt.Sprintf("Error: %v", msg))

	if g.Debug {
		panic(e.String())
	}

	return e
}

type ReadE struct {
	BasicE
	pos Pos
}

func (e *ReadE) Init(g *G, pos Pos, msg string) *ReadE {
	e.BasicE.Init(g, msg)
	e.pos = pos
	return e
}

func (e *ReadE) String() string {
	p := &e.pos

	return fmt.Sprintf(
		"Read error in '%s'; row %v, col %v:\n%v",
		p.src, p.Row, p.Col, e.msg)
}

func (g *G) ReadE(pos Pos, msg string, args ...interface{}) *ReadE {
	msg = fmt.Sprintf(msg, args...)
	e := new(ReadE).Init(g, pos, msg)
	return e
}
