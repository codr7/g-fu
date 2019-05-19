package gfu

import (
  //"log"
  "math"
)

func div(g *G, task *Task, env *Env, args Vec) (Val, E) {
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

func mod(g *G, task *Task, env *Env, args Vec) (Val, E) {
  rhs, e := g.Int(args[1])

  if e != nil {
    return nil, e
  }
  
  return Int(args[0].(Int) % rhs), nil
}

func (e *Env) InitMath(g *G) {
  e.AddFun(g, "div", div, A("x"), A("y"))
  e.AddFun(g, "mod", mod, A("x"), A("y"))

  var pi Dec
  pi.SetFloat(math.Pi)
  e.AddConst(g, "PI", pi)
}
