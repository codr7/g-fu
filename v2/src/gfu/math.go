package gfu

import (
	//"log"
	"math"
	"math/rand"
)

func int_div_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
	lhs, e := g.Int(args[0])

	if e != nil {
		return nil, e
	}

	rhs, e := g.Int(args[1])

	if e != nil {
		return nil, e
	}

	return lhs / rhs, nil
}

func mod_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
	rhs, e := g.Int(args[1])

	if e != nil {
		return nil, e
	}

	return Int(args[0].(Int) % rhs), nil
}

func rand_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
	a := args[0]

	if a == &g.NIL {
		return Int(rand.Int()), nil
	}

	max, e := g.Int(a)

	if e != nil {
		return nil, e
	}

	return Int(rand.Intn(int(max))), nil
}

func (e *Env) InitMath(g *G) {
	e.AddPun(g, "div", int_div_imp, A("x"), A("y"))
	e.AddPun(g, "mod", mod_imp, A("x"), A("y"))
	e.AddFun(g, "rand", rand_imp, AOpt("max", nil))

	var pi Float
	pi.SetFloat(math.Pi)
	e.AddConst(g, "PI", pi)
}
