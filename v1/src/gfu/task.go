package gfu

import (
  "log"
)

type Task struct {
  Inbox Chan
  body Vec

  recall bool
  recall_args Vec
  result Val
}

func (t *Task) Init(inbox Chan, body Vec) *Task {
  t.Inbox = inbox
  t.body = body
  return t
}

func (t *Task) Start(g *G, root_env *Env) {
  var env Env
  root_env.Clone(&env)
  
  go func () {
    var e E
    
    if t.result, e = t.body.EvalExpr(g, t, &env); e != nil {
      log.Fatal(e)
    }
  }()
}
