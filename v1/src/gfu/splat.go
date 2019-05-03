package gfu

import (
	//"log"
	"strings"
)

type Splat struct {
	Wrap
}

func NewSplat(g *G, val Val) *Splat {
	s := new(Splat)
	s.Wrap.Init(&g.SplatType, s, val)
	return s
}

func (s *Splat) Dump(out *strings.Builder) {
	s.val.Dump(out)
	out.WriteString("..")
}

func (s *Splat) Eq(g *G, rhs Val) bool {
	rs, ok := rhs.(*Splat)

	if !ok {
		return false
	}

	return s.val.Eq(g, rs.val)
}

func (s *Splat) Eval(g *G, task *Task, env *Env) (v Val, e E) {
	if v, e = s.val.Eval(g, task, env); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (s *Splat) Expand(g *G, task *Task, env *Env, depth Int) (v Val, e E) {
	if v, e = s.val.Expand(g, task, env, depth); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (s *Splat) Quote(g *G, task *Task, env *Env) (v Val, e E) {
	if v, e = s.val.Quote(g, task, env); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (s *Splat) Splat(g *G, out Vec) (Vec, E) {
	v := s.val

	switch v := v.(type) {
	case Vec:
		return v.Splat(g, out)
	case *Nil:
		return out, nil
	default:
		break
	}

	return append(out, s), nil
}
