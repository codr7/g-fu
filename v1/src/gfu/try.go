package gfu

type Try struct {
  restarts Env
}

func (t *Try) Init(prev *Try) *Try {
  if prev != nil {
    prev.restarts.Dup(&t.restarts)
  }

  return t
}

func (t *Try) AddRestart(g *G, id *Sym, imp *Fun) E {
  return t.restarts.Let(g, id, t.NewRestart(id, imp))
}
