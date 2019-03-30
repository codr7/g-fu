package gfu

type BoolType struct {
  BasicType
}

func (t *BoolType) Init(id *Sym) *BoolType {
  t.BasicType.Init(id)
  return t
}

func (t *BoolType) AsBool(g *G, val Val) bool {
  return val.imp.(bool)
}

func (v Val) AsBool(g *G) bool {
  vt := v.val_type
  
  if vt == g.Bool {
    return v.imp.(bool)
  }

  return vt.AsBool(g, v)
}
