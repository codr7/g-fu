package gfu

import (
  //"log"
  "strings"
)

type MacroType struct {
  BasicType
}

func (t *MacroType) Call(g *G, pos Pos, val Val, args []Form, env *Env) (Val, E) {
  m := val.AsMacro()
  var f Form
  var e E  
  avs := make([]Val, len(args))
  
  for i, a := range args {
    avs[i], e = a.Quote(g, env, 1)
  }
  
  if f, e = m.Call(g, pos, avs, env); e != nil {
    return g.NIL, e
  }
  
  return f.Eval(g, env)
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
