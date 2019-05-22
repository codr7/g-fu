package gfu

import (
  "bufio"
  //"log"
  "os"
)

func flush_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  w := (*bufio.Writer)(args[0].(*Writer))
  
  if e := w.Flush(); e != nil {
    return nil, g.E("Failed flushing writer: %v", e)
  }
  
  return &g.NIL, nil
}

func print_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  out := (*bufio.Writer)(args[0].(*Writer))

  for _, v := range args[1:] {
    g.Print(v, out)
  }

  return &g.NIL, nil
}

func (e *Env) InitIO(g *G) {
  e.AddFun(g, "flush", flush_imp, A("out"))
  e.AddFun(g, "print", print_imp, A("out"), ASplat("vals"))

  e.AddConst(g, "CR", Byte(13))
  e.AddConst(g, "LF", Byte(10))
  
  e.Let(g, g.Sym("stderr"), NewWriter(os.Stderr))
  e.Let(g, g.Sym("stdout"), NewWriter(os.Stdout))
}
