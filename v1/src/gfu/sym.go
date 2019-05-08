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

func (s *Sym) Lookup(g *G, task *Task, env *Env) (v Val, _ *Env, e E) {
  max := len(s.parts)

  for i, p := range s.parts {
    if v, e = env.Get(g, task, p); e != nil {
      return nil, nil, e
    }

    if i == max-1 {
      break
    }
    
    var ok bool
    
    if env, ok = v.(*Env); !ok {
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
  s := val.(*Sym)
  switch s {
  case g.this_sym:
    return env, nil
  default:
    break
  }
  
  v, _, e = s.Lookup(g, task, env)
  return v, e
}

func (_ *SymType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  s := val.(*Sym).parts[0]
  return dst.Extend(g, src, clone, s)
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
