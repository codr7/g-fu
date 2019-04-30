package gfu

import (
  //"log"
)

func div(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return Int(args[0].(Int) / args[1].(Int)), nil
}

func mod(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return Int(args[0].(Int) % args[1].(Int)), nil
}

func (e *Env) InitMath(g *G) {
  e.AddFun(g, "div", div, A("x"), A("y"))
  e.AddFun(g, "mod", mod, A("x"), A("y"))
}
