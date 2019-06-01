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

func (t *Try) AddRestart(g *G, imp *Fun) E {
  return t.restarts.Let(g, imp.id, t.NewRestart(imp))
}
