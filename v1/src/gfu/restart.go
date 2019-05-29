package gfu

type Try struct {
  restarts []*Sym
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
