package gfu

import (
  "fmt"
  "log"
  "strings"
  "sync"
)

type Task struct {
  BasicVal
  Inbox Chan
  
  id           *Sym
  body         Vec
  mutex        sync.Mutex
  cond         *sync.Cond
  done,
  safe         bool
  result       Val
}

func NewTask(g *G, env *Env, id *Sym, inbox Chan, safe bool, body Vec) (*Task, E) {
  return new(Task).Init(g, env, id, inbox, safe, body)
}

func (t *Task) Init(g *G, env *Env, id *Sym, inbox Chan, safe bool, body Vec) (*Task, E) {
  t.BasicVal.Init(&g.TaskType, t)

  if id != nil {
    t.id = id
    
    if e := env.Let(g, id, t); e != nil {
      return nil, e
    }
  }
  
  t.safe = safe
  t.body = body
  t.cond = sync.NewCond(&t.mutex)
  t.Inbox = inbox
  return t, nil
}

func (t *Task) Bool(g *G) bool {
  t.mutex.Lock()
  out := t.done
  t.mutex.Unlock()
  return out
}

func (t *Task) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "(Task %v)", (chan Val)(t.Inbox))
}

func (t *Task) Start(g *G, env *Env) E {
  var te Env
  
  if e := t.body.Extenv(g, env, &te, t.safe); e != nil {
    return e
  }

  go func() {
    var e E

    if t.result, e = t.body.EvalExpr(g, t, &te); e != nil {
      log.Fatal(e)
    }

    t.mutex.Lock()
    t.done = true
    t.cond.Broadcast()
    t.mutex.Unlock()
  }()

  return nil
}

func (t *Task) Wait() Val {
  t.mutex.Lock()

  for !t.done {
    t.cond.Wait()
  }

  t.mutex.Unlock()
  return t.result
}
