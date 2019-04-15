package gfu

import (
  "fmt"
  //"log"
  "strings"
  "sync/atomic"
)

type Sym struct {
  BasicVal
  tag  Tag
  name string
}

func NewSym(g *G, tag Tag, name string) *Sym {
  return new(Sym).Init(g, tag, name)
}

func (s *Sym) Init(g *G, tag Tag, name string) *Sym {
  s.BasicVal.Init(&g.SymType, s)
  s.tag = tag
  s.name = name
  return s
}

func (s *Sym) Dump(out *strings.Builder) {
  out.WriteString(s.name)
}

func (s *Sym) Eval(g *G, task *Task, env *Env) (Val, E) {
  _, found := env.Find(s)

  if found == nil {
    return nil, g.E("Unknown: %v", s)
  }

  return found.Val, nil
}

func (s *Sym) String() string {
  return s.name
}

func (g *G) GSym(prefix string) *Sym {
  var name string
  tag := g.NextSymTag()

  if len(prefix) > 0 {
    name = fmt.Sprintf("g-%v-%v", prefix, tag)
  } else {
    name = fmt.Sprintf("g-%v", tag)
  }

  s := NewSym(g, tag, name)
  g.syms.Store(name, s)
  return s
}

func (g *G) Sym(name string) *Sym {
  var s Sym

  if out, found := g.syms.LoadOrStore(name, &s); found {
    return out.(*Sym)
  }

  return s.Init(g, g.NextSymTag(), name)
}

func (g *G) NextSymTag() Tag {
  return Tag(atomic.AddUint64(&g.nsyms, 1))
}
