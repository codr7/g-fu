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
  nargs := len(args)
  
  if (m.min_args != -1 && nargs < m.min_args) ||
    (m.max_args != -1 && nargs > m.max_args) {
    return g.NIL, g.E(pos, "Arg mismatch")
  }

  var v Val
  var e E
  
  if m.imp == nil {
    var be Env
    m.env.Clone(&be)
    
    for i, a := range m.args {
      id := a.name
      
      if strings.HasSuffix(id, "..") {
        v := new(Vec)
        v.items = make([]Val, nargs-i)

        for _, va := range args[i:] {
          q, e := va.Quote(g, env, 1)

          if e != nil {
            return g.NIL, e
          }

          v.Push(q)
        }
        
        var vv Val
        vv.Init(g.Vec, v)
        be.Put(g.S(id[:len(id)-2]), vv)
        break
      }

      q, e := args[i].Quote(g, env, 1)

      if e != nil {
        return g.NIL, e
      }
      
      be.Put(a, q)
    }
    
    if v, e = Forms(m.body).Eval(g, &be); e != nil {
      return g.NIL, e
    }
  } else {
    var f Form
    
    if f, e = m.imp(g, pos, args, env); e != nil {
      return g.NIL, e
    }

    v, e = f.Quote(g, env, 1)
    
    if e != nil {
      return g.NIL, e
    }
  }

  return v, nil
}

func (t *MacroType) Dump(val Val, out *strings.Builder) {
  m := val.AsMacro()
  out.WriteString("(macro (")

  for i, a := range m.args {
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
