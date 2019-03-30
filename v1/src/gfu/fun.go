package gfu

import (
  "strings"
)

type Fun struct {
  min_args, max_args int
  args []*Sym
  body []Form
  env *Env
}

func NewFun(args []*Sym, body []Form, env *Env) *Fun {
  return new(Fun).Init(args, body, env)
}

func (f *Fun) Init(args []*Sym, body []Form, env *Env) *Fun {
  f.args = args
  f.body = body
  f.env = env
  
  nargs := len(args)

  if nargs > 0 {
    f.min_args, f.max_args = nargs, nargs
    a := args[nargs-1]
    
    if strings.HasSuffix(a.name, "..") {
      f.min_args--
      f.max_args = -1
    }
  }
  
  return f
}
