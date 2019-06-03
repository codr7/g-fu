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

func (g *G) Try(task *Task, env, args_env *Env, body func() (Val, E), restarts...*Fun) (v Val, e E){
  prev := task.try
  var t Try 
  task.try = t.Init(prev)  
  defer func() { task.try = prev }()
  t.AddRestart(g, g.abort_fun)
  t.AddRestart(g, g.retry_fun)
  
  for _, rf := range restarts {
    if e = t.AddRestart(g, rf); e != nil {
      return nil, e
    }
  }
restart:
  v, e = body()
  var ok bool
  
  if e != nil {
    if _, ok = e.(Abort); ok {
      return nil, e
    }

    if _, ok = e.(Retry); ok {
      goto restart
    }

    var rv Val
    var ce E
    
    if rv, ce = g.Catch(task, env, e, args_env); ce != nil {
      if _, ok = ce.(Abort); ok {
        return nil, e
      }
      
      if _, ok = ce.(Retry); ok {
        goto restart
      }
      
      e = ce
    }

    if rv == nil {
      v, e = g.BreakLoop(task, env, e, args_env)
    } else {
      var r Restart
      
      if r, ok = rv.(Restart); !ok {
        return nil, g.E("Expected Restart: %v", rv.Type(g))
      }
      
      if r.try == &t {
        return r.imp.CallArgs(g, task, env, r.args, args_env)
      } else {
        e = r
      }
    }

    if e != nil {
      if _, ok = e.(Abort); ok {
        return nil, e
      }
      
      if _, ok = e.(Retry); ok {
        goto restart
      }
    }
  }

  return v, e
}
