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

func (t *Try) End(g *G) E {
  for _, id := range t.restarts {
    if t.task.restarts.Remove(id) == nil {
      return g.E("Failed removing restart: %v", id)
    }
  }

  return nil
}

func (t *Task) AddRestart(g *G, id *Sym, f *Fun) E {
  if !t.restarts.Add(id, f) {
    return g.E("Dup restart: %v", id)
  }

  try := t.try

  if try == nil {
    return g.E("Restart outside of try")    
  }

  try.restarts = append(try.restarts, id)
  return nil
}
