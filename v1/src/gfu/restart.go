package gfu

import (
  "bufio"
)

type Restart struct {
  try *Try
  imp *Fun
}

type RestartType struct {
  BasicType
}

func (t *Try) NewRestart(imp *Fun) (r Restart) {
  r.try = t
  r.imp = imp
  return r
}

func (r Restart) String() string {
  return "Restart"
}

func (_ Restart) Type(g *G) Type {
  return &g.RestartType
}

func (_ *RestartType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (v Val, e E) {
  r := val.(Restart)
  
  if v, e = g.Call(task, env, r.imp, args, args_env); e != nil {
    return nil, e
  }

  return nil, r
}

func (_ *RestartType) Dump(g *G, val Val, out *bufio.Writer) E {
  return g.Dump(val.(Restart).imp, out)
}

func (t *Task) AddRestart(g *G, id *Sym, f *Fun) E {
  try := t.try

  if try == nil {
    return g.E("Restart outside of try")    
  }

  if !t.restarts.Add(id, try.NewRestart(f)) {
    return g.E("Dup restart: %v", id)
  }

  try.restarts = append(try.restarts, id)
  return nil
}
