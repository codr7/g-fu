package gfu

import (
  "strings"
)

type VecType struct {
  BasicType
}

func (t *VecType) Init(id *Sym) *VecType {
  t.BasicType.Init(id)
  return t
}

func (t *VecType) AsBool(g *G, val Val) bool {
  return val.AsVec().items != nil
}

func (t *VecType) Dump(val Val, out *strings.Builder) {
  v := val.AsVec()
  out.WriteRune('(')
  
  for i, iv := range v.items {
    if i > 0 {
      out.WriteRune(' ')
    }

    iv.Dump(out)
  }
  
  out.WriteRune(')')
}

func (t *VecType) New(g *G, val Val, args ListForm, env *Env, pos Pos) (Val, Error)  {
  is, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, g.NewError(pos, "Constructor arg eval failed: %v", e)
  }

  var out Val
  v := new(Vec)
  v.items = is
  out.Init(g.Vec, v)
  return out, nil
}

func (t *VecType) Splat(g *G, val Val, out []Val) []Val {
  return append(out, val.AsVec().items...)
}

func (v Val) AsVec() *Vec {
  return v.imp.(*Vec)
}

