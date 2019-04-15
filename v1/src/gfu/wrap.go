package gfu

import (
  //"log"
)

type Wrap struct {
  BasicVal
  val Val
}

func (w *Wrap) Init(imp_type *Type, imp Val, val Val) *Wrap {
  w.BasicVal.Init(imp_type, imp)
  w.val = val
  return w
}
  
func (w Wrap) Bool(g *G) bool {
  return w.val.Bool(g)
}

func (w Wrap) Clone(g *G) (Val, E) {
  var e E
  
  if w.val, e = w.val.Clone(g); e != nil {
    return nil, e
  }
  
  return w.imp, nil
}

func (w Wrap) Dup(g *G) (Val, E) {
  var e E
  
  if w.val, e = w.val.Dup(g); e != nil {
    return nil, e
  }
  
  return w.imp, nil
}

func (w Wrap) Pop(g *G) (Val, Val, E) {
  return nil, nil, g.E("Pop not supported: %v", w.imp_type)
}

func (w Wrap) Push(g *G, its...Val) (Val, E) {
  return nil, g.E("Push not supported: %v", w.imp_type)
}
