package gfu

import (
  "fmt"
  "strings"
)

type PrimImp func (*G, *Task, *Env, Vec) (Val, E)

type Prim struct {
  id *Sym
  arg_list ArgList
  imp PrimImp
}

func NewPrim(g *G, id *Sym, imp PrimImp, args []Arg) *Prim {
  p := new(Prim)
  p.id = id
  p.arg_list.Init(g, args)
  p.imp = imp
  return p
}

func (_ *Prim) Bool(g *G) bool {
  return true
}

func (p *Prim) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  if e := p.arg_list.Check(g, args); e != nil {
    return nil, e
  }

  return p.imp(g, task, env, args)
}

func (p *Prim) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "(Prim %v)", p.id)
}

func (p *Prim) Eq(g *G, rhs Val) bool {
  return p == rhs
}

func (p *Prim) Eval(g *G, task *Task, env *Env) (Val, E) {
  return p, nil
}

func (p *Prim) Is(g *G, rhs Val) bool {
  return p == rhs
}

func (p *Prim) Quote(g *G, task *Task, env *Env) (Val, E) {
  return p, nil
}

func (p *Prim) Splat(g *G, out Vec) Vec {
  return append(out, p)
}

func (p *Prim) Type(g *G) *Type {
  return &g.PrimType
}

func (env *Env) AddPrim(g *G, id string, imp PrimImp, args...string) E {
  ids := g.Sym(id)
  vs := make(Vec, len(args))

  for i, a := range args {
    vs[i] = g.Sym(a)
  }

  as, e := ParseArgs(g, vs)

  if e != nil {
    return e
  }

  env.Put(ids, NewPrim(g, ids, imp, as))
  return nil
}
