package gfu

import (
  //"log"
  "strings"
)

type Val struct {
  pos Pos
  val_type Type
  imp interface{}
}

func (v *Val) Init(pos Pos, val_type Type, imp interface{}) *Val {
  v.pos = pos
  v.val_type = val_type
  v.imp = imp
  return v
}
  
func (v Val) Call(g *G, pos Pos, args []Val, env *Env) (Val, E) {
  return v.val_type.Call(g, pos, v, args, env)
}

func (v Val) Dump(out *strings.Builder) {
  v.val_type.Dump(v, out)
}

func (v Val) Eq(g *G, rhs Val) bool {
  return v.val_type.Eq(g, v, rhs)
}

func (v Val) Eval(g *G, pos Pos, env *Env) (Val, E) {
  return v.val_type.Eval(g, pos, v, env)
}

func (v Val) Is(g *G, rhs Val) bool {
  return v.val_type.Is(g, v, rhs)
}

func (v Val) Quote(g *G, pos Pos, env *Env) (Val, E) {
  return v.val_type.Quote(g, pos, v, env)
}

func (v Val) Splat(g *G, pos Pos, out []Val) []Val {
  return v.val_type.Splat(g, pos, v, out)
}

func (v Val) String() string {
  return DumpString(v)
}

type List []Val

func (in List) Eval(g *G, pos Pos, env *Env) ([]Val, E) {
  var out []Val
  
  for _, iv := range in {
    ov, e := iv.Eval(g, pos, env)

    if e != nil {
      return nil, g.E(iv.pos, "Arg eval failed: %v", e)
    }

    if g.recall {
      break
    }
    
    if ov.val_type == g.SplatType {
      out = ov.Splat(g, pos, out)
    } else {
      if ov.val_type == g.VecType {
        v := ov.AsVec()
        v.items = ov.Splat(g, pos, nil)
      }
      
      out = append(out, ov)
    }
  }

  return out, nil
}

type Expr []Val

func (vs Expr) Eval(g *G, pos Pos, env *Env) (Val, E) {
  out := g.NIL
  
  for _, v := range vs {
    var e E
    
    if out, e = v.Eval(g, pos, env); e != nil {
      return g.NIL, e
    }

    if g.recall {
      break
    }
  }

  return out, nil
}

func (env *Env) AddVal(g *G, id string, val_type Type, val interface{}, out *Val) {
  out.Init(NIL_POS, val_type, val)
  env.Put(g.Sym(id), *out)
}
