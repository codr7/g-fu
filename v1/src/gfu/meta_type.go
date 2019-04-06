package gfu

import (
  //"log"
  "strings"
)

type MetaType struct {
  BasicType
}

func (t *MetaType) Bool(g *G, val Val) bool {
  return false
}

func (t *MetaType) Call(g *G, pos Pos, val Val, args []Form, env *Env) (Val, E) {
  vs, e := VecForm(args).Eval(g, env)

  if e != nil {
    return g.NIL, e
  }
  
  return val.AsMeta().New(g, pos, val, vs, env)
}
  
func (t *MetaType) Dump(val Val, out *strings.Builder) {
  out.WriteString(val.AsMeta().Id().name)
}

func (v Val) AsMeta() Type {
  return v.imp.(Type)
}
