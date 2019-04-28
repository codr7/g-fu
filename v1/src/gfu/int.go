package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Int int64

type IntIter struct {
  BasicIter
  pos, max Int
}

func (i Int) Abs() Int {
  if i < 0 {
    return -i
  }

  return i
}

func (i Int) Bool(g *G) bool {
  return i != 0
}

func (i Int) Call(g *G, task *Task, env *Env, args Vec) (Val, E) {
  return nil, g.E("Call not supported: Int")
}

func (i Int) Clone(g *G) (Val, E) {
  return i, nil
}

func (_ Int) Drop(g *G, n Int) (Val, E) {
  return nil, g.E("Drop not supported: Int")
}

func (i Int) Dup(g *G) (Val, E) {
  return i, nil
}

func (i Int) Dump(out *strings.Builder) {
  fmt.Fprintf(out, "%v", int64(i))
}

func (i Int) Eq(g *G, rhs Val) bool {
  return i.Is(g, rhs)
}

func (i Int) Eval(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Expand(g *G, task *Task, env *Env, depth Int) (Val, E) {
  return i, nil
}

func (i Int) Extenv(g *G, src, dst *Env, clone bool) E {
  return nil
}

func (i Int) Is(g *G, rhs Val) bool {
  return i == rhs
}

func (i Int) Iter(g *G) (Val, E) {
  return new(IntIter).Init(g, i), nil
}

func (_ Int) Len(g *G) (Int, E) {
  return -1, g.E("Len not supported: Int")
}

func (_ Int) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: Int")
}

func (i Int) Print(out *strings.Builder) {
  i.Dump(out)
}

func (_ Int) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: Int")
}

func (i Int) Quote(g *G, task *Task, env *Env) (Val, E) {
  return i, nil
}

func (i Int) Splat(g *G, out Vec) (Vec, E) {
  return append(out, i), nil
}

func (i Int) String() string {
  return DumpString(i)
}

func (_ Int) Type(g *G) *Type {
  return &g.IntType
}

func (i *IntIter) Init(g *G, max Int) *IntIter {
  i.BasicVal.Init(&g.IterType, i)
  i.max = max
  return i
}

func (i *IntIter) Bool(g *G) bool {
  return i.pos < i.max
}

func (i *IntIter) Drop(g *G, n Int) (Val, E) {
  if i.max - i.pos < n {
    return nil, g.E("Nothing to drop")
  }

  i.pos += n
  return i, nil
}

func (i *IntIter) Dup(g *G) (Val, E) {
  out := *i
  return &out, nil
}

func (i *IntIter) Eq(g *G, rhs Val) bool {
  ri, ok := rhs.(*IntIter)
  return ok && ri.max == i.max && ri.pos == i.pos
}

func (i *IntIter) Pop(g *G) (Val, Val, E) {
  if i.pos >= i.max {
    return nil, nil, g.E("Nothing to pop")
  }

  v := i.pos
  i.pos++
  return v, i, nil
}

func (i *IntIter) Splat(g *G, out Vec) (Vec, E) {
  for {
    v, _, e := i.Pop(g)
    
    if e != nil {
      return nil, e
    }

    if v == nil {
      break
    }
    
    out = append(out, v)
  }

  return out, nil
}
