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
  
  body         Vec
  mutex        sync.Mutex
  cond         *sync.Cond
  done,
  recall,
  safe         bool
  recall_args  Vec
  result       Val
}

func NewTask(g *G, inbox Chan, safe bool, body Vec) *Task {
  return new(Task).Init(g, inbox, safe, body)
}

func (t *Task) Init(g *G, inbox Chan, safe bool, body Vec) *Task {
  t.BasicVal.Init(&g.TaskType, t)
  t.Inbox = inbox
  t.safe = safe
  t.body = body
  t.cond = sync.NewCond(&t.mutex)
  return t
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

func (t *Task) Start(g *G, env *Env) {
  var te Env
  
  if t.safe {
    env.Clone(g, &te)
  } else {
    env.Dup(g, &te)
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
}

func (t *Task) Wait() Val {
  t.mutex.Lock()

  for !t.done {
    t.cond.Wait()
  }

  t.mutex.Unlock()
  return t.result
}
