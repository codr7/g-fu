package gfu

import (
	"fmt"
	"strings"
)

type Chan chan Val

type ChanType struct {
	BasicType
}

func NewChan(buf Int) Chan {
	return make(Chan, buf)
}

func (c Chan) Dump(out *strings.Builder) {
	fmt.Fprintf(out, "(Chan %v)", (chan Val)(c))
}

func (c Chan) Len() Int {
	return Int(len(c))
}

func (c Chan) Pop() Val {
	return <-c
}

func (c Chan) Push(g *G, its ...Val) E {
	for _, v := range its {
		var e E

		if v, e = g.Clone(v); e != nil {
			return e
		}

		c <- v
	}

	return nil
}

func (_ Chan) Type(g *G) Type {
	return &g.ChanType
}

func (_ *ChanType) Bool(g *G, val Val) (bool, E) {
	return len(val.(Chan)) != 0, nil
}

func (_ *ChanType) Drop(g *G, val Val, n Int) (Val, E) {
	c := val.(Chan)

	for i := Int(0); i < n; i++ {
		<-c
	}

	return c, nil
}

func (_ *ChanType) Dump(g *G, val Val, out *strings.Builder) E {
	val.(Chan).Dump(out)
	return nil
}

func (_ *ChanType) Len(g *G, val Val) (Int, E) {
	return val.(Chan).Len(), nil
}

func (_ *ChanType) Pop(g *G, val Val) (Val, Val, E) {
	v := val.(Chan).Pop()

	if v == nil {
		v = &g.NIL
	}

	return v, val, nil
}

func (_ *ChanType) Push(g *G, val Val, its ...Val) (Val, E) {
	if e := val.(Chan).Push(g, its...); e != nil {
		return nil, e
	}

	return val, nil
}
