package gfu

import (
	//"log"
	"strings"
)

type Splat struct {
	BasicWrap
}

type SplatType struct {
	BasicWrapType
}

func NewSplat(g *G, val Val) (s Splat) {
	s.BasicWrap.Init(val)
	return s
}

func (_ Splat) Type(g *G) Type {
	return &g.SplatType
}

func (_ *SplatType) Dump(g *G, val Val, out *strings.Builder) E {
	if e := g.Dump(val.(Splat).val, out); e != nil {
		return e
	}

	out.WriteString("..")
	return nil
}

func (_ *SplatType) Eval(g *G, task *Task, env *Env, val Val) (v Val, e E) {
	if v, e = g.Eval(task, env, val.(Splat).val); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (_ *SplatType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (v Val, e E) {
	if v, e = g.Expand(task, env, val.(Splat).val, depth); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (_ *SplatType) Quote(g *G, task *Task, env *Env, val Val) (v Val, e E) {
	if v, e = g.Quote(task, env, val.(Splat).val); e != nil {
		return nil, e
	}

	return NewSplat(g, v), nil
}

func (_ *SplatType) Splat(g *G, val Val, out Vec) (Vec, E) {
	s := val.(Splat)
	v := s.val

	switch v := v.(type) {
	case Vec:
		return g.Splat(v, out)
	case *Nil:
		return out, nil
	default:
		break
	}

	return append(out, s), nil
}

func (_ *SplatType) Unwrap(val Val) Val {
	return val.(Splat).val
}

func (_ *SplatType) Wrap(g *G, val Val) Val {
	return NewSplat(g, val)
}
