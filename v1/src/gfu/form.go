package gfu

import (
  "fmt"
  "log"
  "strings"
)

type Form interface {
  fmt.Stringer
  Body() []Form
  Eval(g *G, env *Env) (*Val, Error)
  FormType() *FormType
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
}

func (f *BasicForm) Init(form_type *FormType) *BasicForm {
  f.form_type = form_type
  return f
}

func (f *BasicForm) Body() []Form {
  return []Form{f}
}

func (f *BasicForm) Eval(g *G, env *Env) (*Val, Error) {
  return nil, nil
}

func (f *BasicForm) FormType() *FormType {
  return f.form_type
}

func (f *BasicForm) String() string {
  return f.form_type.id
}

type ExprForm struct {
  BasicForm
  body []Form
}

func (f *ExprForm) Init() *ExprForm {
  f.BasicForm.Init(&FORM_EXPR)
  return f
}

func (f *ExprForm) Append(forms...Form) {
  f.body = append(f.body, forms...)
}

func (f *ExprForm) Body() []Form {
  return f.body
}

func (f *ExprForm) Eval(g *G, env *Env) (*Val, Error) {
  b := f.body
  
  if len(b) > 0 {
    bf := b[0]
    
    if (bf.FormType() == &FORM_ID) {
      switch bid := bf.(*IdForm).id; bid {
      case g.Sym("fun"):
        asf := b[1]

        if asf.FormType() != &FORM_EXPR {
          return nil, g.NewError(&g.Pos, "Invalid fun args: %v", asf)
        }

        var as []*Sym
        
        for _, af := range asf.(*ExprForm).body {
          if af.FormType() != &FORM_ID {
            return nil, g.NewError(&g.Pos, "Invalid fun arg: %v", af)
          }

          as = append(as, af.(*IdForm).id)
        }
        
        return NewVal(g.Fun, NewFun(as, b[2:], env)), nil
      default:
        log.Printf("call %v", bid)
      }
    }
  }

  return nil, nil
}

func (f *ExprForm) String() string {
  var buf strings.Builder
  buf.WriteRune('(')

  for i, bf := range f.body {
    if i > 0 {
      buf.WriteRune(' ')
    }
    
    buf.WriteString(bf.String())
  }
  
  buf.WriteRune(')')
  return buf.String()
}

type IdForm struct {
  BasicForm
  id *Sym
}

func (f *IdForm) Init(id *Sym) *IdForm {
  f.BasicForm.Init(&FORM_ID)
  f.id = id
  return f
}

func (f *IdForm) String() string {
  return f.id.name
}

type LitForm struct {
  BasicForm
  val Val
}

func (f *LitForm) Init(val *Val) *LitForm {
  f.BasicForm.Init(&FORM_LIT)
  f.val = *val
  return f
}


func (f *LitForm) String() string {
  return f.val.String()
}
