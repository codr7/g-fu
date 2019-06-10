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

type OpsType struct {
	BasicType
}

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

func (_ Ops) Type(g *G) Type {
	return &g.OpsType
}

func (_ *OpsType) Dump(g *G, val Val, out *bufio.Writer) E {
	out.WriteRune('(')
	
	for i, op := range val.(Ops) {
		if i > 0 {
			out.WriteRune(' ')
		}
		
		out.WriteString(op.OpId(g).name)
	}

	out.WriteRune(')')
	return nil
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
	form Val
	env Env
	body Ops
}

func NewLetOp(form Val) *LetOp {
	return &LetOp{form: form}
}

func (op *LetOp) DumpArgs(g *G, out *bufio.Writer, depth int) E {
	out.WriteRune(' ')
	op.env.Dump(g, out)
	return nil
}

func (op *LetOp) Eval(g *G, task *Task, env, args_env *Env) (val Val, e E) {
	var let_env *Env

	if len(op.body) == 0 {
		let_env = env
	} else {
		let_env = new(Env)
		
		if e = g.Extenv(env, let_env, op.form, false); e != nil {
			return nil, e
		}
	}
	
	if e = g.Extenv(args_env, let_env, op.form, false); e != nil {
		return nil, e
	}
	
	val = &g.NIL
	
	for _, v := range op.env.vars {
		if val, e = v.Val.(Ops).Eval(g, task, env, args_env); e != nil {
			return nil, e
		}
		
		if e = env.Let(g, v.key, val); e != nil {
			return nil, e
		}
	}

	
	if len(op.body) == 0 {
		return val, nil
	}

	if val, e = op.body.Eval(g, task, let_env, args_env); e != nil {
		return nil, e
	}

	return val, nil

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
