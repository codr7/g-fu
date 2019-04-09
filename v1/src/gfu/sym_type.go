package gfu

import (
  "fmt"
  "strings"
)

type SymType struct {
  BasicType
}

func (t *SymType) Dump(val Val, out *strings.Builder) {
  out.WriteString(val.AsSym().name)
}

func (t *SymType) Eval(g *G, pos Pos, val Val, env *Env) (Val, E) {
  s := val.AsSym()
  _, found := env.Find(s)

  if found == nil {
    return g.NIL, g.E(val.pos, "Unknown: %v", s)
  }

  return found.Val, nil
}

func (t *SymType) New(g *G, pos Pos, val Val, args []Val, env *Env) (Val, E)  {
  n := fmt.Sprintf("g%v", g.NextSymTag())
  
  var out Val
  out.Init(pos, g.SymType, g.Sym(n))
  return out, nil
}

func (v Val) AsSym() *Sym {
  return v.imp.(*Sym)
}
