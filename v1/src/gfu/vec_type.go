package gfu

import (
  "strings"
)

type VecType struct {
  BasicType
}

func (t *VecType) Bool(g *G, val Val) bool {
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

func (t *VecType) Eq(g *G, x Val, y Val) bool {
  xv, yv := x.AsVec().items, y.AsVec().items

  if len(xv) != len(yv) {
    return false
  }

  for i, xi := range xv {
    yi := yv[i]
    
    if !xi.Eq(g, yi) {
      return false
    }
  }

  return true
}

func (t *VecType) New(g *G, pos Pos, val Val, args VecForm, env *Env) (Val, E)  {
  is, e := args.Eval(g, env)

  if e != nil {
    return g.NIL, g.E(pos, "Constructor arg eval failed: %v", e)
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

func (t *VecType) Unquote(g *G, pos Pos, val Val) (Form, E) {
  f := new(ExprForm).Init(pos)

  for _, v := range val.AsVec().items {
    vf, e := v.Unquote(g, pos)

    if e != nil {
      return nil, e
    }

    f.Append(vf)
  }
  
  return f, nil
}

func (v Val) AsVec() *Vec {
  return v.imp.(*Vec)
}

