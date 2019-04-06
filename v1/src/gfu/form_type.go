package gfu

import (
  //"log"
  "strings"
)

type FormType struct {
  BasicType
}

func (t *FormType) Dump(val Val, out *strings.Builder) {
  val.AsForm().Dump(out)
}

func (t *FormType) Unquote(g *G, pos Pos, val Val) (Form, E) {
  return val.AsForm(), nil
}

func (v Val) AsForm() Form {
  return v.imp.(Form)
}
