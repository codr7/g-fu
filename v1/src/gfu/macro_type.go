package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type MacroType struct {
  BasicType
}

func (t *MacroType) Call(g *G, pos Pos, val Val, args []Form, env *Env) (Val, E) {
  m := val.AsMacro()
  var e E
  
  if m.imp == nil {
    avs := make([]Val, len(args))
    
    for i, a := range args {
      avs[i], e = a.Quote(g, env, 1)
    }
    
    return m.CallBody(g, pos, avs, env)
  }

  var f Form
  
  if f, e = m.CallImp(g, pos, args, env); e != nil {
    return g.NIL, e
  }
  
  return f.Quote(g, env, 1)
}

func (t *MacroType) Dump(val Val, out *strings.Builder) {
  m := val.AsMacro()
  out.WriteString("(macro (")

  for i, a := range m.arg_list.args {
    if i > 0 {
      out.WriteRune(' ')
    }

    out.WriteString(a.name)
  }

  if m.imp == nil {
    fmt.Fprintf(out, ") %v)", m.imp)
  } else {
    out.WriteString(") ")
    
    for i, bf := range m.body {
      if i > 0 {
        out.WriteRune(' ')
      }
      
      out.WriteString(bf.String())   
    }
  
    out.WriteRune(')')
  }
}

func (v Val) AsMacro() *Macro {
  return v.imp.(*Macro)
}
