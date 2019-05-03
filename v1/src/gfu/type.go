package gfu

import (
	"strings"
)

type Type struct {
	BasicVal
	id *Sym
}

func (t *Type) Init(g *G, id *Sym) *Type {
	t.BasicVal.Init(&g.MetaType, t)
	t.id = id
	return t
}

func (t *Type) Dump(out *strings.Builder) {
	out.WriteString(t.id.name)
}

func (t *Type) Id() *Sym {
	return t.id
}

func (e *Env) AddType(g *G, t *Type, id string) E {
	t.Init(g, g.Sym(id))
	return e.Let(g, t.Id(), t)
}
