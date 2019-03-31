package gfu

import (
  "fmt"
  "strings"
)

type Form interface {
  fmt.Stringer
  Dumper
  
  Body() []Form
  Eval(g *G, env *Env) (Val, Error)
  FormType() *FormType
  Pos() Pos
}

type FormType struct {
  id string
}

var FORM_EXPR, FORM_ID, FORM_LIT FormType

func init() {
  FORM_EXPR.Init("Expr")
  FORM_ID.Init("Id")
  FORM_LIT.Init("Lit")
}

func (t *FormType) Init(id string) *FormType {
  t.id = id
  return t
}

type BasicForm struct {
  form_type *FormType
  pos Pos
}

func (f *BasicForm) Init(form_type *FormType, pos Pos) *BasicForm {
  f.form_type = form_type
  f.pos = pos
  return f
}

func (f *BasicForm) Body() []Form {
  return []Form{f}
}

func (f *BasicForm) Dump(out *strings.Builder) {
  out.WriteString(f.form_type.id)
}

func (f *BasicForm) Eval(g *G, env *Env) (Val, Error) {
  return g.NIL, nil
}

func (f *BasicForm) FormType() *FormType {
  return f.form_type
}

func (f *BasicForm) Pos() Pos {
  return f.pos
}

func (f *BasicForm) String() string {
  return DumpString(f)
}

type ExprForm struct {
  BasicForm
  body []Form
}

func (f *ExprForm) Init(pos Pos) *ExprForm {
  f.BasicForm.Init(&FORM_EXPR, pos)
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
  
  if b == nil {
    return g.NIL, nil
  }
  
  bf := b[0]
  fv, e := bf.Eval(g, env)
  
  if e != nil {
    return g.NIL, e
  }
  
  rv, e := fv.Call(g, b[1:], env, bf.Pos())
  
  if e != nil {
    return g.NIL, g.E(bf.Pos(), "Call failed: %v", e)
  }
  
  return rv, nil
}

func (f *ExprForm) String() string {
  return DumpString(f)
}

type IdForm struct {
  BasicForm
  id *Sym
}

func (f *IdForm) Init(pos Pos, id *Sym) *IdForm {
  f.BasicForm.Init(&FORM_ID, pos)
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

func (f *IdForm) String() string {
  return DumpString(f)
}

type ListForm []Form

func (f ListForm) Eval(g *G, env *Env) ([]Val, Error) {
  var out []Val
  
  for _, bf := range f {
    v, e := bf.Eval(g, env)

    if e != nil {
      return nil, g.E(bf.Pos(), "Arg eval failed: %v", e)
    }

    if v.val_type == g.Splat {
      out = v.AsSplat().Splat(g, out)
    } else {
      out = append(out, v)
    }
  }

  return out, nil
}

type LitForm struct {
  BasicForm
  val Val
}

func (f *LitForm) Init(pos Pos, val Val) *LitForm {
  f.BasicForm.Init(&FORM_LIT, pos)
  f.val = val
  return f
}

func (f *LitForm) Eval(g *G, env *Env) (Val, Error) {
  return f.val, nil
}

func (f *LitForm) Dump(out *strings.Builder) {
  f.val.Dump(out)
}

func (f *LitForm) String() string {
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
  }

  return out, nil

}
