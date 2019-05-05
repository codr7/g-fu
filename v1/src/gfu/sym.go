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
  return s
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

func (_ *SymType) Eval(g *G, task *Task, env *Env, val Val) (Val, E) {
  return env.Get(g, val.(*Sym))
}

func (_ *SymType) Extenv(g *G, src, dst *Env, val Val, clone bool) E {
  s := val.(*Sym)
  
  if i, dv := dst.Find(s); dv == nil {
    if _, sv := src.Find(s); sv != nil {
      if clone {
        dv = dst.Insert(i, sv.key)
        var e E

        if dv.Val, e = g.Clone(sv.Val); e != nil {
          return e
        }
      } else {
        dst.InsertVar(i, sv)
      }
    }
  }

  return nil
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
