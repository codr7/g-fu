package gfu

import (
  "bufio"
  //"log"
  "os"
)

func print_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  out := (*bufio.Writer)(args[0].(*Writer))

  for _, v := range args[1:] {
    g.Print(v, out)
  }

  return &g.NIL, nil
}

func (e *Env) InitIO(g *G) {
  e.AddFun(g, "print", print_imp, A("out"), ASplat("vals"))

  e.Let(g, g.Sym("stderr"), NewWriter(os.Stderr))
  e.Let(g, g.Sym("stdout"), NewWriter(os.Stdout))
}
