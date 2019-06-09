package gfu

import (
	"bufio"
	"bytes"
	"fmt"
	//"log"
)

type Bin []byte

type BinType struct {
	BasicType
}

type BinIter struct {
	in  Bin
	pos Int
}

type BinIterType struct {
	BasicIterType
}

func NewBin(len int) Bin {
	return make(Bin, len)
}

func (b Bin) Type(g *G) Type {
	return &g.BinType
}

func (_ *BinType) Bool(g *G, val Val) (bool, E) {
	return len(val.(Bin)) > 0, nil
}

func (_ *BinType) Dump(g *G, val Val, out *bufio.Writer) E {
	out.WriteString("(0x")

	for _, v := range val.(Bin) {
		fmt.Fprintf(out, " %02x", v)
	}

	out.WriteRune(')')
	return nil
}

func (_ *BinType) Dup(g *G, val Val) (Val, E) {
	var dst Bin
	src := val.(Bin)

	if len(src) > 0 {
		dst = NewBin(len(src))
		copy(dst, src)
	}

	return dst, nil
}

func (_ *BinType) Eq(g *G, lhs, rhs Val) (bool, E) {
	return bytes.Compare(lhs.(Bin), rhs.(Bin)) == 0, nil
}

func (_ *BinType) Index(g *G, val Val, key Vec) (Val, E) {
	if len(key) > 1 {
		return nil, g.E("Invalid index: %v", key.Type(g))
	}

	b := val.(Bin)
	i, ok := key[0].(Int)

	if !ok {
		return nil, g.E("Invalid index: %v", key[0].Type(g))
	}

	if i := int(i); i < 0 || i > len(b) {
		return nil, g.E("Index out of bounds: %v", i)
	}

	return Byte(b[i]), nil
}

func (_ *BinType) Iter(g *G, val Val) (Val, E) {
	return new(BinIter).Init(g, val.(Bin)), nil
}

func (_ *BinType) Len(g *G, val Val) (Int, E) {
	return Int(len(val.(Bin))), nil
}

func (_ *BinType) Print(g *G, val Val, out *bufio.Writer) E {
	out.WriteString(string(val.(Bin)))
	return nil
}

func (_ *BinType) SetIndex(g *G, task *Task, env *Env, val Val, key Vec, set Setter) (Val, Val, E) {
	if len(key) > 1 {
		return nil, nil, g.E("Invalid index: %v", key.Type(g))
	}

	b := val.(Bin)
	i, ok := key[0].(Int)

	if !ok {
		return nil, nil, g.E("Invalid index: %v", key[0].Type(g))
	}

	if i := int(i); i < 0 || i > len(b) {
		return nil, nil, g.E("Index out of bounds: %v", i)
	}

	v, e := set(Byte(b[i]))

	if e != nil {
		return nil, nil, e
	}

	bv, ok := v.(Byte)

	if !ok {
		return nil, nil, g.E("Expected Byte: %v", v.Type(g))
	}

	b[i] = byte(bv)
	return bv, b, nil
}

func (i *BinIter) Init(g *G, in Bin) *BinIter {
	i.in = in
	return i
}

func (_ *BinIter) Type(g *G) Type {
	return &g.BinIterType
}

func (_ *BinIterType) Bool(g *G, val Val) (bool, E) {
	i := val.(*BinIter)
	return i.pos < Int(len(i.in)), nil
}

func (_ *BinIterType) Drop(g *G, val Val, n Int) (Val, E) {
	i := val.(*BinIter)

	if Int(len(i.in))-i.pos < n {
		return nil, g.E("Nothing to drop")
	}

	i.pos += n
	return i, nil
}

func (_ *BinIterType) Dup(g *G, val Val) (Val, E) {
	out := *val.(*BinIter)
	return &out, nil
}

func (_ *BinIterType) Eq(g *G, lhs, rhs Val) (bool, E) {
	li := lhs.(*BinIter)
	ri, ok := rhs.(*BinIter)

	if !ok {
		return false, nil
	}

	ok, e := g.Eq(ri.in, li.in)

	if e != nil {
		return false, e
	}

	return ok && ri.pos == li.pos, nil
}

func (_ *BinIterType) Pop(g *G, val Val) (Val, Val, E) {
	i := val.(*BinIter)

	if i.pos >= Int(len(i.in)) {
		return &g.NIL, i, nil
	}

	v := Byte(i.in[i.pos])
	i.pos++
	return v, i, nil
}

func (_ *BinIterType) Splat(g *G, val Val, out Vec) (Vec, E) {
	i := val.(*BinIter)

	for _, v := range i.in[i.pos:] {
		out = append(out, Byte(v))
	}

	return out, nil
}
