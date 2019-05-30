package gfu

type Try struct {
  task *Task
  restarts []*Sym
}

func (t *Task) NewTry() *Try {
  try := new(Try)
  try.task = t
  return try
}

func (t *Try) AddRestart(g *G, id *Sym, f *Fun) E {
  if !t.task.restarts.Add(id, f) {
    return g.E("Dup restart: %v", id)
  }

  t.restarts = append(t.restarts, id)
  return nil
}

func (t *Try) End(g *G) E {
  for _, id := range t.restarts {
    if t.task.restarts.Remove(id) == nil {
      return g.E("Failed removing restart: %v", id)
    }
  }

  t.restarts = nil
  return nil
}
