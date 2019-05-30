package gfu

import (
  "bufio"
  "fmt"
  "os"
  "strings"
)

type Abort struct {}

func (_ Abort) String() string {
  return "Abort"
}

type Retry struct {}

func (r Retry) String() string {
  return "Retry"
}

func (g *G) BreakLoop(task *Task, env *Env, cause E, args_env *Env) (Val, E) {
  fmt.Printf("%v\n", cause)
  rs := task.restarts.vars
  stdin := bufio.NewReader(os.Stdin) 
  stdout := bufio.NewWriter(os.Stdout)
  
  for i, v := range rs {
    fmt.Printf("%v %v", i, v.key)
    v.Val.(*Fun).arg_list.items.Dump(g, stdout)
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
  r := rs[n].Val.(*Fun)
  return g.Call(task, env, r, args, args_env)
}
