package gfu

import (
	"bufio"
	"fmt"
	//"log"
	"strings"
	"sync/atomic"
)

type Sym struct {
	tag   Tag
	name  string
	parts []*Sym
}

type SymType struct {
	BasicType
}

func NewSym(g *G, tag Tag, name string) *Sym {
	return new(Sym).Init(g, tag, name)
}

func (s *Sym) Init(g *G, tag Tag, name string) *Sym {
	s.tag = tag
	s.name = name

	if strings.IndexRune(name, '/') != -1 {
		for _, p := range strings.Split(name, "/") {
			if len(p) > 0 {
				s.parts = append(s.parts, g.Sym(p))
			}
		}
	}

	if s.parts == nil {
		s.parts = append(s.parts, s)
	}

	return s
}

func (s *Sym) LookupVar(g *G, env *Env, silent bool) (v *Var, i int, _ *Env, args []Val, e E) {
	max := len(s.parts)

	for j, p := range s.parts {
		if v, i, e = env.GetVar(g, p, silent); e != nil {
			return nil, 0, env, nil, e
		}

		if (silent && v == nil) || j == max-1 {
			break
		}

		var ok bool
		t := v.Val.Type(g)

		if t == &g.MetaType {
			t = v.Val.(Type)
		} else {
			args = append(args, NewQuote(g, v.Val))
		}

		if env, ok = v.Val.(*Env); !ok {
			env = t.Env()
		}
	}

	if v == nil {
		sps := s.parts
		s := sps[len(sps)-1]
		sn := s.name

		if sn[0] == '$' {
			i, _ = env.Find(s)
			v = env.Insert(i, s)
			v.Val = g.NewSym(sn[1:])
			return v, i, env, args, nil
		}
	}

	return v, i, env, args, nil
}

func (s *Sym) Lookup(g *G, task *Task, env, args_env *Env, silent bool) (Val, *Env, []Val, E) {
	var v *Var
	var val Val
	var args []Val

	if v, _, env, args, _ = s.LookupVar(g, env, true); v == nil {
		val, _ = env.Resolve(g, task, s.parts[len(s.parts)-1], args_env, true)

		if val == nil && !silent {
			return nil, env, args, g.EUnknown(s)
		}
	} else {
		val = v.Val
	}

	return val, env, args, nil
}

func (s *Sym) String() string {
	return s.name
}

func (s *Sym) Suffix() *Sym {
	ps := s.parts
	return ps[len(ps)-1]
}

func (_ *Sym) Type(g *G) Type {
	return &g.SymType
}

func (_ *SymType) Dump(g *G, val Val, out *bufio.Writer) E {
	out.WriteString(val.(*Sym).name)
	return nil
}

func (_ *SymType) Eval(g *G, task *Task, env *Env, val Val, args_env *Env) (v Val, e E) {
	s := val.(*Sym)

	use_val := NewFun(g, env, g.Sym("use-val"), A("new"))
	use_val.imp = func(g *G, task *Task, env *Env, args Vec) (Val, E) {
		return args[0], nil
	}

	v, e = g.Try(task, env, args_env, func() (Val, E) {
		if v, args_env, _, _ = s.Lookup(g, task, env, args_env, true); v == nil {

			return nil, g.EUnknown(s)
		}

		return v, nil
	}, use_val)

	if e != nil {
		return nil, e
	}

	if p, ok := v.(*Prim); ok && p.arg_list.items == nil {
		v, e = g.Call(task, env, v, nil, args_env)
	} else if m, ok := v.(*Mac); ok && m.arg_list.items == nil {
		v, e = g.Call(task, env, v, nil, args_env)
	}

	if task.pure > 0 && args_env != env {
		v, e = g.Clone(v)
	}

	return v, e
}

func (_ *SymType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (v Val, e E) {
	s := val.(*Sym)

	if v, _, _, e = s.Lookup(g, task, env, env, true); e != nil {
		return nil, e
	}

	if v != nil {
		if m, ok := v.(*Mac); ok {
			if m.arg_list.items == nil {
				return m.ExpandCall(g, task, env, nil)
			}
		}
	}

	return val, nil
}

func (_ *SymType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
	s := val.(*Sym).parts[0]

	if s.name[0] == '$' {
		return nil
	}

	return dst.Extend(g, src, clone, s)
}

func (g *G) NewSym(prefix string) *Sym {
	var name string
	tag := g.NextSymTag()

	if len(prefix) > 0 {
		name = fmt.Sprintf("%v-sym-%v", prefix, tag)
	} else {
		name = fmt.Sprintf("sym-%v", tag)
	}

	s := NewSym(g, tag, name)
	g.syms.Store(name, s)
	return s
}

func (g *G) Sym(name string, args ...interface{}) *Sym {
	var s Sym

	if len(args) > 0 {
		name = fmt.Sprintf(name, args...)
	}

	if out, found := g.syms.LoadOrStore(name, &s); found {
		return out.(*Sym)
	}

	return s.Init(g, g.NextSymTag(), name)
}

func (g *G) NextSymTag() Tag {
	return Tag(atomic.AddUint64(&g.nsyms, 1))
}
