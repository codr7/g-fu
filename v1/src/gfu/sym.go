package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Sym struct {
  tag Tag
  name string
}

func NewSym(g *G, name string) *Sym {
  return new(Sym).Init(g, name)
}

func (s *Sym) Init(g *G, name string) *Sym {
  s.tag = Tag(len(g.syms))
  s.name = name
  g.syms[name] = s
  return s
}

func (s *Sym) Bool(g *G) bool {
  return true
}

func (s *Sym) Call(g *G, args Vec, env *Env) (Val, E) {
  return s, nil 
}

func (s *Sym) Dump(out *strings.Builder) {
  out.WriteString(s.name)
}

func (s *Sym) Eq(g *G, rhs Val) bool {
  return s == rhs
}

func (s *Sym) Eval(g *G, env *Env) (Val, E) {
  _, found := env.Find(s)

  if found == nil {
    return nil, g.E("Unknown: %v", s)
  }

  return found.Val, nil
}

func (s *Sym) Is(g *G, rhs Val) bool {
  return s == rhs
}

func (s *Sym) Quote(g *G, env *Env) (Val, E) {
  return s, nil
}

func (s *Sym) Splat(g *G, out Vec) Vec {
  return append(out, s)
}

func (s *Sym) String() string {
  return s.name
}

func (_ *Sym) Type(g *G) *Type {
  return &g.SymType
}

func (g *G) GSym(prefix string) *Sym {
  var name string
  n := len(g.syms)
    
  if len(prefix) > 0 {
    name = fmt.Sprintf("g-%v-%v", prefix, n)
  } else {
    name = fmt.Sprintf("g-%v", n)
  }
  
  return g.Sym(name)
}

func (g *G) Sym(name string) *Sym {
  if s := g.syms[name]; s != nil { return s }
  return new(Sym).Init(g, name)
}

