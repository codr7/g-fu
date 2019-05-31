package gfu

type Catch struct {
  etype Type
  imp *Fun
}

func (c *Catch) Init(etype Type, imp *Fun) *Catch {
  c.etype = etype
  c.imp = imp
  return c
}

func (g *G) Catch(task *Task, env *Env, val Val, args_env *Env) (Val, E) {
  t := val.Type(g)
  
  for _, c := range task.catch_q {
    if vt.Isa(c.etype) {
      return g.Call(task, env, c.imp, Vec{val}, args_env)
    }
  }

  return nil, nil
}

