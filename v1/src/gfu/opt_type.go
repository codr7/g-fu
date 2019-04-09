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

func (v Val) AsOpt() Val {
  return v.imp.(Val)
}
