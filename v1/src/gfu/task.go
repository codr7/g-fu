package gfu

import (
  "fmt"
  "log"
  "strings"
  "sync"
)

type Task struct {
  Inbox Chan
  
  body Vec
  mutex sync.Mutex
  exit *sync.Cond
  done, recall bool
  recall_args Vec
  result Val
}

func NewTask(g *G, inbox Chan, body Vec) *Task {
  return new(Task).Init(g, inbox, body)
}

func (t *Task) Init(g *G, inbox Chan, body Vec) *Task {
  t.Inbox = inbox
  t.body = body
  t.exit = sync.NewCond(&t.mutex)
  return t
}

func (t *Task) Bool(g *G) (out bool) {
  t.mutex.Lock()
  out = t.done
  t.mutex.Unlock()
  return out
}

func (t *Task) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return t, nil
}

func (t *Task) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "(Task %v)", t)
}

func (t *Task) Eq(g *G, rhs Val) bool {
  return t.Is(g, rhs)
}

func (t *Task) Eval(g *G, task *Task, env *Env) (Val, E) {
  return t, nil
}

func (t *Task) Is(g *G, rhs Val) bool {
  return t == rhs
}

func (t *Task) Quote(g *G, task *Task, env *Env) (Val, E) {
  return t, nil
}

func (t *Task) Splat(g *G, out Vec) Vec {
  return append(out, t)
}

func (t *Task) Start(g *G, root_env *Env) {
  var env Env
  root_env.Clone(&env)
  
  go func () {
    var e E
    
    if t.result, e = t.body.EvalExpr(g, t, &env); e != nil {
      log.Fatal(e)
    }

    t.mutex.Lock()
    t.done = true
    t.exit.Broadcast()
    t.mutex.Unlock()
  }()
}

func (t *Task) Type(g *G) *Type {
  return &g.TaskType
}

func (t *Task) Wait() Val {
  t.mutex.Lock()

  if !t.done {
    t.exit.Wait()
  }
  
  t.mutex.Unlock()
  return t.result
}
