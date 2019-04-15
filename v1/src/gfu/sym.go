package gfu

import (
	"fmt"
	//"log"
	"strings"
	"sync/atomic"
)

type Sym struct {
	tag  Tag
	name string
}

func NewSym(tag Tag, name string) *Sym {
	return new(Sym).Init(tag, name)
}

func (s *Sym) Init(tag Tag, name string) *Sym {
	s.tag = tag
	s.name = name
	return s
}

func (s *Sym) Bool(g *G) bool {
	return true
}

func (s *Sym) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
	return s, nil
}

func (s *Sym) Dump(out *strings.Builder) {
	out.WriteString(s.name)
}

func (s *Sym) Eq(g *G, rhs Val) bool {
	return s == rhs
}

func (s *Sym) Eval(g *G, task *Task, env *Env) (Val, E) {
	_, found := env.Find(s)

	if found == nil {
		return nil, g.E("Unknown: %v", s)
	}

	return found.Val, nil
}

func (s *Sym) Is(g *G, rhs Val) bool {
	return s == rhs
}

func (s *Sym) Quote(g *G, task *Task, env *Env) (Val, E) {
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
	tag := g.NextSymTag()

	if len(prefix) > 0 {
		name = fmt.Sprintf("g-%v-%v", prefix, tag)
	} else {
		name = fmt.Sprintf("g-%v", tag)
	}

	s := NewSym(tag, name)
	g.syms.Store(name, s)
	return s
}

func (g *G) Sym(name string) *Sym {
	var s Sym

	if out, found := g.syms.LoadOrStore(name, &s); found {
		return out.(*Sym)
	}

	return s.Init(g.NextSymTag(), name)
}

func (g *G) NextSymTag() Tag {
	return Tag(atomic.AddUint64(&g.nsyms, 1))
}
