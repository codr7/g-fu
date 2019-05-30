package gfu

type Throw struct {
  val Val
}

func (t Throw) String() string {
  return "Throw"
}

func (g *G) Throw(val Val) E {
  return Throw{val}
}
