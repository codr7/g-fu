package gfu

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

type Abort struct { }

type AbortType struct {
  BasicType
}

type Retry struct { }

type RetryType struct {
  BasicType
}

type Restart struct {
  try *Try
  imp *Fun
  args Vec
}

type RestartType struct {
  BasicType
}

func (_ Abort) Type(g *G) Type {
  return &g.AbortType
}

func (_ *AbortType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("Abort")
  return nil
}

func (_ Retry) Type(g *G) Type {
  return &g.RetryType
}

func (_ RetryType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("Retry")
  return nil
}

func (t *Try) NewRestart(imp *Fun) (r Restart) {
  r.try = t
  r.imp = imp
  return r
}

func (r Restart) Dump(g *G, out *bufio.Writer) E {
  fmt.Fprintf(out, "Restart: %v", g.EString(r.imp))
  return nil
}

func (_ Restart) Type(g *G) Type {
  return &g.RestartType
}

func (_ *RestartType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (Val, E) {
  r := val.(Restart)
  return g.Call(task, env, r.imp, append(r.args, args...), args_env)
}

func (_ *RestartType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString("(restart ")
  if e := g.Dump(val.(Restart).imp, out); e != nil {
    return e
  }
  
  out.WriteRune(')')
  return nil
}

func (g *G) BreakLoop(task *Task, env *Env, cause E, args_env *Env) (Val, E) {
  stdin := bufio.NewReader(os.Stdin) 
  stdout := bufio.NewWriter(os.Stdout)
  rs := task.try.restarts.vars

  stdout.WriteString("\nBreak: ")
  g.Print(cause, stdout)
  stdout.WriteRune('\n')
  stdout.Flush()
  
  for i, v := range rs {
    fmt.Printf("%v %v", i, v.key)
    v.Val.(Restart).imp.arg_list.items.Dump(g, stdout)
    stdout.WriteRune('\n')
    stdout.Flush()
  }

  var in string
  var n Int
  var args Vec
  
  for {
    fmt.Printf("\nChoose 0-%v: ", len(rs))
    in, _ = stdin.ReadString('\n')
    out, e := g.ReadAll(&INIT_POS, strings.NewReader(in), nil)

    if e != nil {
      fmt.Printf("Failed reading choice: %v\n", e)
      continue
    }

    out, e = out.EvalVec(g, task, env, args_env)

    if e != nil {
      fmt.Printf("Failed evaluating choice: %v\n", e)
      continue
    }

    var ok bool
    n, ok = out[0].(Int)

    if !ok {
      fmt.Printf("Expected Int: %v", out[0].Type(g))
      continue      
    }
    
    if n < 0 || int(n) > len(rs)-1 {
      continue
    }

    args = out[1:]
    break
  }
  
  fmt.Printf("\n")
  r := rs[n].Val.(Restart).imp
  return g.Call(task, env, r, args, args_env)
}
