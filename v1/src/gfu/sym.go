package gfu

import (
  "fmt"
  //"log"
  "strings"
  "sync/atomic"
)

type Sym struct {
  tag  Tag
  name string
  parts []*Sym
  root bool
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
        s.parts  = append(s.parts, g.Sym(p))
      }
    }
  }

  if s.parts == nil {
    s.parts = append(s.parts, s)
  }
  
  return s
}

func (s *Sym) Lookup(g *G, task *Task, env *Env, silent bool) (v Val, _ *Env, e E) {
  max := len(s.parts)

  for i, p := range s.parts {
    if v, e = env.Get(g, task, p, silent); e != nil {
      return nil, nil, e
    }

    if silent && v == nil {
      return nil, nil, nil
    }
    
    if i == max-1 {
      break
    }
    
    var ok bool
    
    if env, ok = v.(*Env); !ok {
      if silent {
        return nil, nil, nil
      }
      
      return nil, nil, g.E("Expected env: %v", v.Type(g))
    }
  }

  return v, env, nil
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
  if v, _, e = val.(*Sym).Lookup(g, task, env, false); e != nil {
    return nil, e
  }
  
  if p, ok := v.(*Prim); ok && p.arg_list.items == nil {
    v, e = g.Call(task, env, v, Vec{})
  } else if m, ok := v.(*Mac); ok && m.arg_list.items == nil {
    v, e = g.Call(task, env, v, Vec{})
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
