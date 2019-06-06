package gfu

import (
	"bufio"
	"fmt"
	"io"
	//"log"
)

type Writer bufio.Writer

type WriterType struct {
	BasicType
}

func NewWriter(in io.Writer) *Writer {
	return (*Writer)(bufio.NewWriter(in))
}

func (_ *Writer) Type(g *G) Type {
	return &g.WriterType
}

func (_ *WriterType) Dump(g *G, val Val, out *bufio.Writer) E {
	fmt.Fprintf(out, "writer-%v", (*bufio.Writer)(val.(*Writer)))
	return nil
}
