package gfu

import (
	"bufio"
	"fmt"
)

type PrimImp func(*G, *Task, *Env, *Env, Vec, Ops) (Ops, E)

type Prim struct {
	id       *Sym
	arg_list ArgList
	pure     bool
	imp      PrimImp
}

type PrimType struct {
	BasicType
}

func NewPrim(g *G, id *Sym, pure bool, imp PrimImp, args []Arg) *Prim {
	p := new(Prim)
	p.id = id
	p.arg_list.Init(g, args)
	p.pure = pure
	p.imp = imp
	return p
}

func (p *Prim) Call(g *G, task *Task, env, args_env *Env, args Vec, out Ops) (Ops, E) {
	return p.imp(g, task, env, args_env, args, out)
}

func (_ *Prim) Type(g *G) Type {
	return &g.PrimType
}

func (_ *PrimType) Dump(g *G, val Val, out *bufio.Writer) E {
	p := val.(*Prim)
	fmt.Fprintf(out, "(prim %v", p.id)
	nargs := len(p.arg_list.items)

	if nargs > 0 {
		out.WriteString(" (")
	}

	for i, a := range p.arg_list.items {
		if i > 0 {
			out.WriteRune(' ')
		}

		if e := a.Dump(g, out); e != nil {
			return e
		}
	}

	if nargs > 0 {
		out.WriteRune(')')
	}

	out.WriteRune(')')
	return nil
}

func (env *Env) AddPrim(g *G, id string, pure bool, imp PrimImp, args ...Arg) E {
	s := g.Sym(id)
	return env.Let(g, s, NewPrim(g, s, pure, imp, args))
}

func ParsePrimArgs(g *G, args Val) Vec {
	if args == &g.NIL {
		return nil
	} else if v, ok := args.(Vec); ok {
		if len(v) == 0 {
			return nil
		}

		return v
	}

	return Vec{args}
}
