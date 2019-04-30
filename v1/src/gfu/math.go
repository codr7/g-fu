package gfu

import (
  //"log"
)

func is_zero(g *G, task *Task, env *Env, args Vec) (Val, E) {
  i, ok := args[0].(Int)
  return g.Bool(ok && i == 0), nil
}

func (e *Env) InitMath(g *G) {
  e.AddFun(g, "z?", is_zero, A("val"))
}
