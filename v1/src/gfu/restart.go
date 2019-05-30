package gfu

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

type Abort struct {}

type Retry struct {}

type Restart struct {
  try *Try
  id *Sym
  imp *Fun
  args []Val
}

type RestartType struct {
  BasicType
}

func (_ Abort) Dump(g *G, out *bufio.Writer) E {
  out.WriteString("Abort")
  return nil
}

func (_ Retry) Dump(g *G, out *bufio.Writer) E {
  out.WriteString("Retry")
  return nil
}

func (t *Try) NewRestart(id *Sym, imp *Fun) (r Restart) {
  r.try = t
  r.id = id
  r.imp = imp
  return r
}

func (r Restart) Dump(g *G, out *bufio.Writer) E {
  fmt.Fprintf(out, "Restart: %v", r.id)
  return nil
}

func (_ Restart) Type(g *G) Type {
  return &g.RestartType
}

func (_ *RestartType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "restart-%v", val.(Restart).id)
  return nil
}

func (g *G) BreakLoop(task *Task, env *Env, cause E, args_env *Env) (Val, E) {
  cs, e := g.DumpString(cause)

  if e != nil {
    return nil, e
  }
  
  fmt.Printf("%v\n", cs)
  rs := task.try.restarts.vars
  stdin := bufio.NewReader(os.Stdin) 
  stdout := bufio.NewWriter(os.Stdout)
  
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
