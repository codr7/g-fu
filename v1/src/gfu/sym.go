package gfu

import (
  "fmt"
  //"log"
  "strings"
  "sync/atomic"
)

type Sym struct {
  BasicVal
  tag  Tag
  name string
}

func NewSym(g *G, tag Tag, name string) *Sym {
  return new(Sym).Init(g, tag, name)
}

func (s *Sym) Init(g *G, tag Tag, name string) *Sym {
  s.BasicVal.Init(&g.SymType, s)
  s.tag = tag
  s.name = name
  return s
}

func (s *Sym) Dump(out *strings.Builder) {
  out.WriteString(s.name)
}

func (s *Sym) Eval(g *G, task *Task, env *Env) (Val, E) {
  return env.Get(g, s)
}

func (s *Sym) Extenv(g *G, src, dst *Env, clone bool) E {
  if i, dv := dst.Find(s); dv == nil {
    if _, sv := src.Find(s); sv != nil {
      if clone {
        dv = dst.Insert(i, sv.key)
        var e E
        
        if dv.Val, e = sv.Val.Clone(g); e != nil {
          return e
        }
      } else {
        dst.InsertVar(i, sv)
      }
    }
  }

  return nil
}

func (s *Sym) Print(out *strings.Builder) {
  out.WriteString(s.name)
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
