package gfu

import (
  "bufio"
  "fmt"
  //"log"
  "strings"
  "sync/atomic"
)

type Sym struct {
  tag    Tag
  name   string
  parts  []*Sym
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

func (s *Sym) LookupVar(g *G, env *Env, args []Val, silent bool) (v *Var, i int, _ *Env, _ []Val, e E) {
  max := len(s.parts)

  for j, p := range s.parts {
    if v, i, e = env.GetVar(g, p, silent); e != nil {
      return nil, 0, nil, nil, e
    }

    if (silent && v == nil) || j == max-1 {
      break
    }

    var ok bool

    if env, ok = v.Val.(*Env); !ok {
      t := v.Val.Type(g)

      if t == &g.MetaType {
        t = v.Val.(Type)
      } else {
        args = append(args, v.Val)
      }

      env = t.Env()
    }
  }

  return v, i, env, args, nil
}

func (s *Sym) Lookup(g *G, task *Task, env, args_env *Env, silent bool) (Val, *Env, []Val, E) {
  var v *Var
  var args []Val
  
  if v, _, env, args, _ = s.LookupVar(g, env, nil, true); v != nil {
    return v.Val, env, args, nil
  }

  if env != nil {
    val, _ := env.Resolve(g, task, s.parts[len(s.parts)-1], args_env, true)

    if val != nil {
      return val, env, args, nil
    }
  }

  if !silent {
    return nil, nil, nil, g.E("Unknown: %v", s)
  }

  return nil, nil, nil, nil
}

func (s *Sym) String() string {
  return s.name
}

func (s *Sym) Suffix() *Sym {
  ps := s.parts
  return ps[len(ps)-1]
}

func (_ *Sym) Type(g *G) Type {
  return &g.SymType
}

func (_ *SymType) Dump(g *G, val Val, out *bufio.Writer) E {
  out.WriteString(val.(*Sym).name)
  return nil
}

func (_ *SymType) Eval(g *G, task *Task, env *Env, val Val, args_env *Env) (v Val, e E) {
  if v, args_env, _, e = val.(*Sym).Lookup(g, task, env, env, false); e != nil {
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

  if v, _, _, e = s.Lookup(g, task, env, env, true); e != nil {
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

func (g *G) Sym(name string, args...interface{}) *Sym {
  var s Sym

  if len(args) > 0 {
    name = fmt.Sprintf(name, args...)
  }
  
  if out, found := g.syms.LoadOrStore(name, &s); found {
    return out.(*Sym)
  }

  return s.Init(g, g.NextSymTag(), name)
}

func (g *G) NextSymTag() Tag {
  return Tag(atomic.AddUint64(&g.nsyms, 1))
}
