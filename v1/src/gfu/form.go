package gfu

import (
  "fmt"
  //"log"
  "strings"
)

type Form interface {
  fmt.Stringer
  Dumper
  
  Body() []Form
  Eval(g *G, env *Env) (Val, Error)
  Quote(g *G, env *Env, depth int) (Val, Error)
  Pos() Pos
}

type BasicForm struct {
  pos Pos
}

func (f *BasicForm) Init(pos Pos) *BasicForm {
  f.pos = pos
  return f
}

func (f *BasicForm) Body() []Form {
  return []Form{f}
}

func (f *BasicForm) Dump(out *strings.Builder) {
  out.WriteRune('?')
}

func (f *BasicForm) Eval(g *G, env *Env) (Val, Error) {
  return g.NIL, nil
}

func (f *BasicForm) Pos() Pos {
  return f.pos
}

func (f *BasicForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  return g.NIL, g.E(f.pos, "Not quote implemented")
}

func (f *BasicForm) String() string {
  return DumpString(f)
}

type ExprForm struct {
  BasicForm
  body []Form
}

func (f *ExprForm) Init(pos Pos) *ExprForm {
  f.BasicForm.Init(pos)
  return f
}

func (f *ExprForm) Append(forms...Form) {
  f.body = append(f.body, forms...)
}

func (f *ExprForm) Body() []Form {
  return f.body
}

func (f *ExprForm) Dump(out *strings.Builder) {
  out.WriteRune('(')

  for i, bf := range f.body {
    if i > 0 {
      out.WriteRune(' ')
    }

    bf.Dump(out)
  }
  
  out.WriteRune(')')
}

func (f *ExprForm) Eval(g *G, env *Env) (Val, Error) {
  b := f.body
  
  if len(b) == 0 {
    return g.NIL, nil
  }
  
  bf := b[0]
  fv, e := bf.Eval(g, env)
  
  if e != nil {
    return g.NIL, e
  }
  
  rv, e := fv.Call(g, bf.Pos(), b[1:], env)
  
  if e != nil {
    return g.NIL, g.E(bf.Pos(), "Call failed: %v", e)
  }
  
  return rv, nil
}

func (f *ExprForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  var out Vec
  
  for _, bf := range f.body {
    v, e := bf.Quote(g, env, depth)

    if e != nil {
      return g.NIL, e
    }
    
    if v.val_type == g.Splat {
      out.items = v.Splat(g, out.items)
    } else {
      out.Push(v)
    }
  }

  var v Val
  v.Init(g.Vec, &out)
  return v, nil
}

func (f *ExprForm) String() string {
  return DumpString(f)
}

type IdForm struct {
  BasicForm
  id *Sym
}

func (f *IdForm) Init(pos Pos, id *Sym) *IdForm {
  f.BasicForm.Init(pos)
  f.id = id
  return f
}

func (f *IdForm) Dump(out *strings.Builder) {
  out.WriteString(f.id.name)
}

func (f *IdForm) Eval(g *G, env *Env) (Val, Error) {
  id := f.id
  splat := false
  
  if strings.HasSuffix(id.name, "..") {
    id = g.S(id.name[:len(id.name)-2])
    splat = true
  }
  
  _, found := env.Find(id)

  if found == nil {
    return g.NIL, g.E(f.pos, "Unknown: %v", id)
  }

  v := found.Val
  
  if splat {
    v.Init(g.Splat, v)
  }
  
  return v, nil
}

func (f *IdForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  var v Val
  v.Init(g.Sym, f.id)
  return v, nil
}

func (f *IdForm) String() string {
  return DumpString(f)
}

type LitForm struct {
  BasicForm
  val Val
}

func (f *LitForm) Init(pos Pos, val Val) *LitForm {
  f.BasicForm.Init(pos)
  f.val = val
  return f
}

func (f *LitForm) Eval(g *G, env *Env) (Val, Error) {
  return f.val, nil
}

func (f *LitForm) Dump(out *strings.Builder) {
  f.val.Dump(out)
}

func (f *LitForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  return f.val, nil
}

func (f *LitForm) String() string {
  return DumpString(f)
}

type QuoteForm struct {
  BasicForm
  form Form
}

func (f *QuoteForm) Init(pos Pos, form Form) *QuoteForm {
  f.BasicForm.Init(pos)
  f.form = form
  return f
}

func (f *QuoteForm) Eval(g *G, env *Env) (Val, Error) {
  return f.form.Quote(g, env, 1)
}

func (f *QuoteForm) Dump(out *strings.Builder) {
  out.WriteRune('\'')
  f.form.Dump(out)
}

func (f *QuoteForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  depth++
  v, e := f.form.Quote(g, env, depth)
  depth--
  return v, e
}

func (f *QuoteForm) String() string {
  return DumpString(f)
}

type SplatForm struct {
  BasicForm
  form Form
}

func (f *SplatForm) Init(pos Pos, form Form) *SplatForm {
  f.BasicForm.Init(pos)
  f.form = form
  return f
}

func (f *SplatForm) Eval(g *G, env *Env) (Val, Error) {
  sv, e := f.form.Eval(g, env)

  if e != nil {
    return g.NIL, e
  }

  var v Val
  v.Init(g.Splat, sv)
  return v, nil
}

func (f *SplatForm) Dump(out *strings.Builder) {
  out.WriteString("..")
}

func (f *SplatForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  v, e := f.form.Quote(g, env, depth)

  if e != nil {
    return g.NIL, e
  }
  
  v.Init(g.Splat, v)
  return v, nil
}

func (f *SplatForm) String() string {
  return DumpString(f)
}

type VecForm []Form

func (f VecForm) Eval(g *G, env *Env) ([]Val, Error) {
  var out []Val
  
  for _, bf := range f {
    v, e := bf.Eval(g, env)

    if e != nil {
      return nil, g.E(bf.Pos(), "Arg eval failed: %v", e)
    }

    if g.recall_args != nil {
      break
    }
    
    if v.val_type == g.Splat {
      out = v.Splat(g, out)
    } else {
      out = append(out, v)
    }
  }

  return out, nil
}

type UnquoteForm struct {
  BasicForm
  form Form
}

func (f *UnquoteForm) Init(pos Pos, form Form) *UnquoteForm {
  f.BasicForm.Init(pos)
  f.form = form
  return f
}

func (f *UnquoteForm) Eval(g *G, env *Env) (Val, Error) {
  return g.NIL, g.E(f.pos, "Nothing to unquote")
}

func (f *UnquoteForm) Dump(out *strings.Builder) {
  out.WriteRune('%')
  f.form.Dump(out)
}

func (f *UnquoteForm) Quote(g *G, env *Env, depth int) (Val, Error) {
  var v Val
  var e Error
  depth--
  
  if depth == 0 {
    v, e = f.form.Eval(g, env)
  } else {
    v, e = f.form.Quote(g, env, depth)
  }
  
  depth++
  return v, e
}

func (f *UnquoteForm) String() string {
  return DumpString(f)
}

type Forms []Form

func (fs Forms) Eval(g *G, env *Env) (Val, Error) {
  out := g.NIL
  
  for _, f := range fs {
    var e Error
    
    if out, e = f.Eval(g, env); e != nil {
      return g.NIL, e
    }

    if g.recall_args != nil {
      break
    }
  }

  return out, nil
}
