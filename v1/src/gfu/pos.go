package gfu

import (
  "fmt"
  "strings"
)

type Pos struct {
  src string
  Row, Col int
}

var (
  INIT_POS, NIL_POS Pos
)

func init() {
  INIT_POS.Init("n/a")
  
  NIL_POS.src = "n/a"
  NIL_POS.Row = -1
  NIL_POS.Col = -1
}

func (p *Pos) Init(src string) *Pos {
  p.src = src
  p.Row = 1
  p.Col = 0
  return p
}

func (p Pos) Bool(g *G) bool {
  return p != NIL_POS && p != INIT_POS;
}

func (p Pos) Call(g *G, args Vec, env *Env) (Val, E) {
  return p, nil 
}

func (p Pos) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "(Pos \"%v\" %v %v)", p.src, p.Row, p.Col)
}

func (p Pos) Eq(g *G, rhs Val) bool {
  rp, ok := rhs.(Pos)
  return ok && p == rp
}

func (p Pos) Eval(g *G, env *Env) (Val, E) {
  return p, nil
}

func (p Pos) Is(g *G, rhs Val) bool {
  return p.Eq(g, rhs)
}

func (p Pos) Quote(g *G, env *Env) (Val, E) {
  return p, nil
}

func (p Pos) Splat(g *G, out Vec) Vec {
  return append(out, p)
}

func (_ Pos) Type(g *G) *Type {
  return &g.PosType
}


