package gfu

import (
  "fmt"
  "strconv"
)

type Abort struct {}

func NewAbort() (a Abort) {
  return a
}

func (_ Abort) String() string {
  return "Abort"
}

type Retry struct {
  try *Try
}

func NewRetry(try *Try) (r Retry) {
  r.try = try
  return r
}

func (r Retry) String() string {
  return "Retry"
}

func (t *Task) AddRestart(g *G, id *Sym, f *Fun) E {
  try := t.try

  if try == nil {
    return g.E("Restart outside of try")    
  }

  if !t.restarts.Add(id, f) {
    return g.E("Dup restart: %v", id)
  }

  try.restarts = append(try.restarts, id)
  return nil
}

func (g *G) BreakLoop(task *Task, env *Env, cause E, args_env *Env) (Val, E) {
  fmt.Printf("%v\n", cause)
  rs := task.restarts.vars
  
  for i, v := range rs {
    fmt.Printf("%v %v\n", i, v.key)
  }

  var in string
  var n int64
  
  for {
    fmt.Printf("\nChoose 0-%v: ", len(rs))
    _, e := fmt.Scanln(&in)
  
    if e != nil {
      fmt.Printf("Failed reading line: %v\n", e)
      continue
    }

    n, e = strconv.ParseInt(in, 10, 64)

    if e != nil {
      fmt.Printf("Failed parsing choice: %v\n", e)
      continue
    }

    if n < 0 || int(n) > len(rs)-1 {
      continue
    }

    break
  }
  
  fmt.Printf("\n")
  r := rs[n].Val.(*Fun)
  return g.Call(task, env, r, nil, args_env)
}
