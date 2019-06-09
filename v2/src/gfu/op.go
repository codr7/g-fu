package gfu

import (
	"bufio"
)

type Op interface {
	DumpArgs(*G, *bufio.Writer, int) E
	Eval(*G, *Task, *Env, *Env) (Val, E)
	EvalVec(*G, *Task, *Env, *Env, Vec) (Vec, E)
	OpId(*G) *Sym
}

type Ops []Op

func (ops Ops) Dump(g *G, out *bufio.Writer, depth int) (e E) {
	for _, op := range ops {
		for i := 0; i < depth; i++ {
			out.WriteString("  ")
		}

		out.WriteString(op.OpId(g).name)

		if e = op.DumpArgs(g, out, depth+1); e != nil {
			return e
		}

		out.WriteRune('\n')
	}

	return nil
}

func (ops Ops) Eval(g *G, task *Task, env, args_env *Env) (v Val, e E) {
	for _, op := range ops {
		if v, e = op.Eval(g, task, env, args_env); e != nil {
			return nil, e
		}
	}

	return v, nil
}

type IfOp struct {
	cond, x, y Ops
}

func NewIfOp(cond, x, y Ops) *IfOp {
	return &IfOp{cond: cond, x: x, y: y}
}

func (op *IfOp) DumpArgs(g *G, out *bufio.Writer, depth int) E {
	out.WriteRune(':')

	if e := op.cond.Dump(g, out, depth); e != nil {
		return e
	}

	if e := op.x.Dump(g, out, depth); e != nil {
		return e
	}

	if op.y != nil {
		if e := op.y.Dump(g, out, depth); e != nil {
			return e
		}
	}

	return nil
}

func (op *IfOp) Eval(g *G, task *Task, env, args_env *Env) (v Val, e E) {
	if v, e = op.cond.Eval(g, task, env, args_env); e != nil {
		return nil, e
	}

	var bv bool

	if bv, e = g.Bool(v); e != nil {
		return nil, e
	}

	if bv {
		return op.x.Eval(g, task, env, args_env)
	}

	if op.y == nil {
		return &g.NIL, nil
	}

	return op.y.Eval(g, task, env, args_env)
}

func (op *IfOp) EvalVec(g *G, task *Task, env, args_env *Env, out Vec) (_ Vec, e E) {
	var v Val

	if v, e = op.Eval(g, task, env, args_env); e != nil {
		return nil, e
	}

	return append(out, v), nil
}

func (op *IfOp) OpId(g *G) *Sym {
	return g.Sym("if")
}

type LetOp struct {
	key *Sym
	val Ops
}

func NewLetOp(key *Sym, val Ops) *LetOp {
	return &LetOp{key: key, val: val}
}

func (op *LetOp) DumpArgs(g *G, out *bufio.Writer, depth int) E {
	out.WriteRune(' ')

	if e := g.Dump(op.key, out); e != nil {
		return e
	}

	out.WriteRune(':')

	if e := op.val.Dump(g, out, depth); e != nil {
		return e
	}

	return nil
}

func (op *LetOp) Eval(g *G, task *Task, env, args_env *Env) (v Val, e E) {
	if v, e = op.val.Eval(g, task, env, args_env); e != nil {
		return nil, e
	}

	if e = env.Let(g, op.key, v); e != nil {
		return nil, e
	}

	return v, nil
}

func (op *LetOp) EvalVec(g *G, task *Task, env, args_env *Env, out Vec) (_ Vec, e E) {
	var v Val

	if v, e = op.Eval(g, task, env, args_env); e != nil {
		return nil, e
	}

	return append(out, v), nil
}

func (op *LetOp) OpId(g *G) *Sym {
	return g.Sym("let")
}

type LitOp struct {
	val Val
}

func NewLitOp(val Val) *LitOp {
	return &LitOp{val: val}
}

func (op *LitOp) DumpArgs(g *G, out *bufio.Writer, depth int) E {
	out.WriteRune(' ')
	return g.Dump(op.val, out)
}

func (op *LitOp) Eval(g *G, task *Task, env *Env, args_env *Env) (Val, E) {
	return op.val, nil
}

func (op *LitOp) EvalVec(g *G, task *Task, env, args_env *Env, out Vec) (Vec, E) {
	return append(out, op.val), nil
}

func (op *LitOp) OpId(g *G) *Sym {
	return g.Sym("lit")
}
