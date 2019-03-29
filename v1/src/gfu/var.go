package gfu

type Var struct {
  key *Sym
  val Val
}

func (v *Var) Init(key *Sym) *Var {
  v.key = key
  return v
}
