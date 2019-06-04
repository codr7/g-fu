package gfu

import (
  "bufio"
  "fmt"
  "log"
  "sync"
)

type Task struct {
  Inbox Chan

  id      *Sym
  body    Vec
  mutex   sync.Mutex
  cond    *sync.Cond
  try     *Try
  catch_q []Catch
  pure    int
  done    bool
  result  Val
}

type TaskType struct {
  BasicType
}

func NewTask(g *G, env *Env, id *Sym, inbox Chan, body Vec) (*Task, E) {
  return new(Task).Init(g, env, id, inbox, body)
}

func (t *Task) Init(g *G, env *Env, id *Sym, inbox Chan, body Vec) (*Task, E) {
  if id != nil {
    t.id = id

    if e := env.Let(g, id, t); e != nil {
      return nil, e
    }
  }

  t.body = body
  t.cond = sync.NewCond(&t.mutex)
  t.Inbox = inbox
  return t, nil
}

func (t *Task) Start(g *G, env *Env) E {
  var te Env

  if e := g.Extenv(env, &te, t.body, true); e != nil {
    return e
  }

  go func() {
    var e E

    if t.result, e = t.body.EvalExpr(g, t, &te, &te); e != nil {
      log.Fatal(e)
    }

    t.mutex.Lock()
    t.done = true
    t.cond.Broadcast()
    t.mutex.Unlock()
  }()

  return nil
}

func (t *Task) Type(g *G) Type {
  return &g.TaskType
}

func (t *Task) Wait() Val {
  t.mutex.Lock()

  for !t.done {
    t.cond.Wait()
  }

  t.mutex.Unlock()
  return t.result
}

func (_ *TaskType) Bool(g *G, val Val) (bool, E) {
  t := val.(*Task)
  t.mutex.Lock()
  out := t.done
  t.mutex.Unlock()
  return out, nil
}

func (_ *TaskType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("(task")

  if t := val.(*Task); t.id != nil {
    fmt.Fprintf(out, " %v", t.id)
  }

  out.WriteString(")")
  return nil
}
