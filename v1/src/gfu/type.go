package gfu

import (
  "bufio"
  "fmt"
  //"log"
  "sync"
)

type Type interface {
  fmt.Stringer
  Val

  Init(*G, *Sym, []Type) E

  ArgList(*G, Val) (*ArgList, E)
  Bool(*G, Val) (bool, E)
  Call(*G, *Task, *Env, Val, Vec, *Env) (Val, E)
  Clone(*G, Val) (Val, E)
  Drop(*G, Val, Int) (Val, E)
  Dump(*G, Val, *bufio.Writer) E
  Dup(*G, Val) (Val, E)
  EachParent(func(Type))
  Env() *Env
  Eq(*G, Val, Val) (bool, E)
  Eval(*G, *Task, *Env, Val, *Env) (Val, E)
  Expand(*G, *Task, *Env, Val, Int) (Val, E)
  Extenv(*G, *Env, *Env, Val, bool) E
  Id() *Sym
  Index(*G, Val, Vec) (Val, E)
  Is(*G, Val, Val) bool
  Isa(Type) Type
  Iter(*G, Val) (Val, E)
  Len(*G, Val) (Int, E)
  Pop(*G, Val) (Val, Val, E)
  Print(*G, Val, *bufio.Writer)
  Push(*G, Val, ...Val) (Val, E)
  Quote(*G, *Task, *Env, Val, *Env) (Val, E)
  SetIndex(*G, Val, Vec, Setter) (Val, E)
  Splat(*G, Val, Vec) (Vec, E)
}

type BasicType struct {
  id      *Sym
  parents sync.Map
  env     Env
}

type MetaType struct {
  BasicType
}

func (t *BasicType) add_parent(key, val Type) {
  key.EachParent(func(k Type) {
    t.add_parent(k, val)
  })

  t.parents.LoadOrStore(key, val)
}

func type_check_imp(g *G, task *Task, env *Env, args Vec) (Val, E) {
  pt, ok := args[0].(Type)

  if !ok {
    return nil, g.E("Expected Type: %v", args[0].Type(g))
  }

  for _, a := range args[1:] {
    pt = g.Isa(a, pt)
  }
  
  if pt == nil {
    return &g.NIL, nil
  }

  return pt, nil
}

func (t *BasicType) Init(g *G, id *Sym, parents []Type) E {
  t.id = id

  for _, p := range parents {
    t.add_parent(p, p)
  }

  return t.env.Use(g, &g.RootEnv, "do", "fun", "mac")
}

func (t *BasicType) ArgList(g *G, _ Val) (*ArgList, E) {
  return nil, g.E("ArgList not supported: %v", t.id)
}

func (_ *BasicType) Bool(g *G, val Val) (bool, E) {
  return true, nil
}

func (t *BasicType) Call(g *G, task *Task, env *Env, val Val, args Vec, args_env *Env) (Val, E) {
  return nil, g.E("Call not supported: %v", t.id)
}

func (_ *BasicType) Clone(g *G, val Val) (Val, E) {
  return g.Dup(val)
}

func (_ *BasicType) Drop(g *G, val Val, n Int) (out Val, e E) {
  for i := Int(0); i < n; i++ {
    if _, out, e = g.Pop(val); e != nil {
      return nil, e
    }
  }

  return out, nil
}

func (t *BasicType) Dump(g *G, val Val, out *bufio.Writer) E {
  fmt.Fprintf(out, "%v", val)
  return nil
}

func (_ *BasicType) Dup(g *G, val Val) (Val, E) {
  return val, nil
}

func (t *BasicType) EachParent(f func(Type)) {
  t.parents.Range(func(key, _ interface{}) bool {
    f(key.(Type))
    return true
  })
}

func (t *BasicType) Env() *Env {
  return &t.env
}

func (_ *BasicType) Eq(g *G, lhs, rhs Val) (bool, E) {
  return g.Is(lhs, rhs), nil
}

func (_ *BasicType) Eval(g *G, task *Task, env *Env, val Val, args_env *Env) (Val, E) {
  return val, nil
}

func (_ *BasicType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (Val, E) {
  return val, nil
}

func (_ *BasicType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  return nil
}

func (t *BasicType) Id() *Sym {
  return t.id
}

func (t *BasicType) Index(g *G, val Val, key Vec) (Val, E) {
  return nil, g.E("Index not supported: %v", t.id)
}

func (_ *BasicType) Is(g *G, lhs, rhs Val) bool {
  return lhs == rhs
}

func (t *BasicType) Isa(parent Type) Type {
  v, ok := t.parents.Load(parent)

  if !ok {
    return nil
  }

  return v.(Type)
}

func (t *BasicType) Iter(g *G, val Val) (Val, E) {
  return nil, g.E("Iter not supported: %v", t.id)
}

func (t *BasicType) Len(g *G, val Val) (Int, E) {
  return -1, g.E("Len not supported: %v", t.id)
}

func (t *BasicType) Pop(g *G, val Val) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: %v", t.id)
}

func (_ *BasicType) Print(g *G, val Val, out *bufio.Writer) {
  g.Dump(val, out)
}

func (t *BasicType) Push(g *G, val Val, its ...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", t.id)
}

func (_ *BasicType) Quote(g *G, task *Task, env *Env, val Val, args_env *Env) (Val, E) {
  return val, nil
}

func (t *BasicType) SetIndex(g *G, val Val, key Vec, set Setter) (Val, E) {
  return nil, g.E("SetIndex not supported: %v", t.id)
}

func (_ *BasicType) Splat(g *G, val Val, out Vec) (Vec, E) {
  return append(out, val), nil
}

func (t *BasicType) String() string {
  return t.id.name
}

func (_ *BasicType) Type(g *G) Type {
  return &g.MetaType
}

func (_ *MetaType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString(val.(Type).Id().name)
  return nil
}

func (e *Env) AddType(g *G, t Type, id string, parents ...Type) E {
  t.Init(g, g.Sym(id), parents)

  t.Env().AddFun(g, "?",
    func (g *G, task *Task, env *Env, args Vec) (Val, E) {
      return type_check_imp(g, task, env, append(Vec{t}, args...))
    },
    A("val0"), ASplat("vals"))

  return e.Let(g, t.Id(), t)
}

func (g *G) Isa(val Val, parent Type) Type {
  vt := val.Type(g)

  if vt == &g.MetaType {
    vt = val.(Type)
  }

  if vt == parent {
    return vt
  }

  return vt.Isa(parent)
}
