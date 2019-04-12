package gfu

import (
  //"log"
  "strings"
)

type MacroType struct {
  BasicType
}

func (t *MacroType) Call(g *G, val Val, args Vec, env *Env) (v Val, e E) {
  m := val.AsMacro()
  avs := make(Vec, len(args))
  
  for i, a := range args {
    if avs[i], e = a.Quote(g, env); e != nil {
      return g.NIL, e
    }
  }

  if v, e = m.Call(g, avs, env); e != nil {
    return g.NIL, e
  }
  
  return v.Eval(g, env)
}

func (t *MacroType) Dump(val Val, out *strings.Builder) {
  m := val.AsMacro()
  out.WriteString("(macro (")

  for i, a := range m.arg_list.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.id.name)
  }

  out.WriteString(") ")
  
  for i, bf := range m.body {
    if i > 0 {
      out.WriteRune(' ')
    }
    
    out.WriteString(bf.String())   
  }
  
  out.WriteRune(')')
}

func (v Val) AsMacro() *Macro {
  return v.imp.(*Macro)
}
