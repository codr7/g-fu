package gfu

import (
  //"log"
  "strings"
)

type OptType struct {
  BasicType
}

func (t *OptType) Bool(g *G, val Val) bool {
  return val.AsOpt().AsBool(g)
}

func (t *OptType) Dump(val Val, out *strings.Builder) {
  val.AsOpt().Dump(out)
  out.WriteRune('?')
}

func (t *OptType) Eq(g *G, x Val, y Val) bool {
  return x.AsOpt().Is(g, y.AsOpt())
}

func (t *OptType) Eval(g *G, pos Pos, val Val, env *Env) (Val, E) {
  var e E
  val.imp, e = val.AsOpt().Eval(g, pos, env)
  return val, e
}

func (t *OptType) Quote(g *G, pos Pos, val Val, env *Env) (Val, E) {
  var e E

  if val, e = val.AsOpt().Quote(g, pos, env); e != nil {
    return g.NIL, e
  }

  val.Init(pos, g.OptType, val)
  return val, nil
}

func (v Val) AsOpt() Val {
  return v.imp.(Val)
}
