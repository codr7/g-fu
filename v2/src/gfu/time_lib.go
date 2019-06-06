package gfu

import (
	"time"
)

func now_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
	return Time(time.Now()), nil
}

func (e *Env) InitTime(g *G) {
	e.AddType(g, &g.NanosType, "Nanos")
	e.AddType(g, &g.TimeType, "Time")

	e.AddFun(g, "now", now_imp)
}
