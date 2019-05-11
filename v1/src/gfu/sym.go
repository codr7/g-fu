package gfu

import (
  "fmt"
  //"log"
  "strings"
  "sync/atomic"
)

type Sym struct {
  tag   Tag
  name  string
  parts []*Sym
  root  bool
}

type SymType struct {
  BasicType
}

func NewSym(g *G, tag Tag, name string) *Sym {
  return new(Sym).Init(g, tag, name)
}

func (s *Sym) Init(g *G, tag Tag, name string) *Sym {
  s.tag = tag
  s.name = name

  if strings.IndexRune(name, '/') != -1 {
    for _, p := range strings.Split(name, "/") {
      if len(p) > 0 {
        s.parts = append(s.parts, g.Sym(p))
      }
    }
  }

  if s.parts == nil {
    s.parts = append(s.parts, s)
  }

  return s
}

func (s *Sym) LookupVar(g *G, env *Env, silent bool) (v *Var, _ *Env, e E) {
  max := len(s.parts)

  for i, p := range s.parts {
    if v, e = env.GetVar(g, p, silent); e != nil {
      return nil, nil, e
    }

    if (silent && v == nil) || i == max-1 {
      break
    }

    var ok bool

    if env, ok = v.Val.(*Env); !ok {
      if silent {
        return nil, nil, nil
      }

      return nil, nil, g.E("Expected env: %v", v.Val.Type(g))
    }
  }

  return v, env, nil
}


func (s *Sym) Lookup(g *G, task *Task, env *Env, silent bool) (Val, *Env, E) {
  var v *Var
  
  if v, env, _ = s.LookupVar(g, env, true); v != nil {
    return v.Val, env, nil
  }
  
  if env != nil {
    val, _ := env.Resolve(g, task, s.parts[len(s.parts)-1], true)

    if val != nil {
      return val, env, nil
    }
  }

  if !silent {
    return nil, nil, g.E("Unknown: %v", s)
  }
    
  return nil, nil, nil
}

func (s *Sym) String() string {
  return s.name
}

func (_ *Sym) Type(g *G) Type {
  return &g.SymType
}

func (_ *SymType) Dump(g *G, val Val, out *strings.Builder) E {
  out.WriteRune('\'')
  out.WriteString(val.(*Sym).name)
  return nil
}

func (_ *SymType) Eval(g *G, task *Task, env *Env, val Val) (v Val, e E) {
  var args_env *Env

  if v, args_env, e = val.(*Sym).Lookup(g, task, env, false); e != nil {
    return nil, e
  }

  if p, ok := v.(*Prim); ok && p.arg_list.items == nil {
    v, e = g.Call(task, env, v, Vec{}, args_env)
  } else if m, ok := v.(*Mac); ok && m.arg_list.items == nil {
    v, e = g.Call(task, env, v, Vec{}, args_env)
  }

  return v, e
}

func (_ *SymType) Expand(g *G, task *Task, env *Env, val Val, depth Int) (v Val, e E) {
  s := val.(*Sym)

  if v, _, e = s.Lookup(g, task, env, true); e != nil {
    return nil, e
  }

  if v != nil {
    if m, ok := v.(*Mac); ok {
      if m.arg_list.items == nil {
        return m.ExpandCall(g, task, env, Vec{})
      }
    }
  }

  return val, nil
}

func (_ *SymType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  return dst.Extend(g, src, clone, val.(*Sym).parts[0])
}

func (s *SymType) Print(g *G, val Val, out *strings.Builder) {
  out.WriteString(val.(*Sym).name)
}

func (g *G) NewSym(prefix string) *Sym {
  var name string
  tag := g.NextSymTag()

  if len(prefix) > 0 {
    name = fmt.Sprintf("%v-%v", prefix, tag)
  } else {
    name = fmt.Sprintf("sym-%v", tag)
  }

  s := NewSym(g, tag, name)
  g.syms.Store(name, s)
  return s
}

func (g *G) Sym(name string) *Sym {
  var s Sym

  if out, found := g.syms.LoadOrStore(name, &s); found {
    return out.(*Sym)
  }

  return s.Init(g, g.NextSymTag(), name)
}

func (g *G) NextSymTag() Tag {
  return Tag(atomic.AddUint64(&g.nsyms, 1))
}
