package gfu

import (
  //"log"
)

type Catch struct {
  etype Type
  imp *Fun
}

func (c *Catch) Init(etype Type, imp *Fun) *Catch {
  c.etype = etype
  c.imp = imp
  return c
}

func (g *G) Catch(task *Task, env *Env, val Val, args_env *Env) (bool, E) {
  t := val.Type(g)
  
  for _, c := range task.catch_q {
    if c.etype == nil || g.Isa(t, c.etype) != nil {
      _, e := g.Call(task, env, c.imp, Vec{val}, args_env)

      if e == nil {
        return true, val
      }

      return true, e
    }
  }

  return false, val
}

