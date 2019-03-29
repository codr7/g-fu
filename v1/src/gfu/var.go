package gfu

type Var struct {
  key *Sym
  Val Val
}

func (v *Var) Init(key *Sym) *Var {
  v.key = key
  return v
}
