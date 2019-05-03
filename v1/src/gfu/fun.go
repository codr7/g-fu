package gfu

import (
	"fmt"
	//"log"
	"strings"
)

type FunImp func(*G, *Task, *Env, Vec) (Val, E)

type Fun struct {
	BasicVal

	id        *Sym
	env       *Env
	env_cache Env
	arg_list  ArgList
	body      Vec
	imp       FunImp
}

func NewFun(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
	return new(Fun).Init(g, env, id, args)
}

func (f *Fun) Init(g *G, env *Env, id *Sym, args []Arg) (*Fun, E) {
	f.BasicVal.Init(&g.FunType, f)

	if id != nil {
		f.id = id

		if e := env.Let(g, id, f); e != nil {
			return nil, e
		}
	}

	f.env = env
	f.arg_list.Init(g, args)
	return f, nil
}

func (f *Fun) CallArgs(g *G, task *Task, env *Env, args Vec) (Val, E) {
	var e E

	if e = f.arg_list.Check(g, args); e != nil {
		return nil, e
	}

	if f.imp != nil {
		return f.imp(g, task, env, f.arg_list.Fill(g, args))
	}

	var be Env

	if f.env_cache.vars == nil {
		if e = f.body.Extenv(g, f.env, &be, false); e != nil {
			return nil, e
		}

		be.Dup(g, &f.env_cache)
	} else {
		f.env_cache.Dup(g, &be)
	}

recall:
	f.arg_list.LetVars(g, &be, args)
	var v Val

	if v, e = f.body.EvalExpr(g, task, &be); e != nil {
		if r, ok := e.(Recall); ok {
			be.Clear()
			f.env_cache.Dup(g, &be)
			args = r.args
			goto recall
		}

		return nil, e
	}

	return v, e
}

func (f *Fun) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
	args, e := args.EvalVec(g, task, env)

	if e != nil {
		return nil, e
	}

	return f.CallArgs(g, task, env, args)
}

func (f *Fun) Dump(out *strings.Builder) {
	if id := f.id; id == nil {
		out.WriteString("(fun (")
	} else {
		fmt.Fprintf(out, "(fun %v (", f.id)
	}

	for i, a := range f.arg_list.items {
		if i > 0 {
			out.WriteRune(' ')
		}

		out.WriteString(a.String())
	}

	if f.imp == nil {
		out.WriteString(")")

		for _, bv := range f.body {
			out.WriteRune(' ')
			bv.Dump(out)
		}

		out.WriteRune(')')
	} else {
		fmt.Fprintf(out, ") %v)", f.imp)
	}
}

func (env *Env) AddFun(g *G, id string, imp FunImp, args ...Arg) E {
	f, e := NewFun(g, env, g.Sym(id), args)

	if e != nil {
		return e
	}

	f.imp = imp
	return nil
}

type Recall struct {
	args Vec
}

func NewRecall(args Vec) (r Recall) {
	r.args = args
	return r
}

func (r Recall) Dump(out *strings.Builder) {
	out.WriteString("(recall")

	for _, a := range r.args {
		a.Dump(out)
	}

	out.WriteRune(')')
}

func (r Recall) String() string {
	return DumpString(r)
}
