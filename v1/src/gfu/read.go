package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func unread(g *G, in *strings.Reader) Error {
  if e := in.UnreadRune(); e != nil {
    return g.NewError(g.Pos, "Error unreading char")
  }

  return nil
}

func (g *G) Read(in *strings.Reader, end rune) (Form, Error) {
  var c rune
  var e Error

  for {
    var re error
    c, _, re = in.ReadRune()

    if re == io.EOF {
      return nil, nil
    }
    
    if re != nil {
      return nil, g.NewError(g.Pos, "Error reading char: %v", re)
    }
    
    if c == end {
      return nil, nil
    }
    
    switch c {
    case ' ':
      g.Pos.Col++
      continue
    case '\n':
      g.Pos.Col = MIN_POS.Col
      g.Pos.Row++
      continue
    case '(':
      return g.ReadExpr(in)
    default:
      if unicode.IsDigit(c) {
        if e = unread(g, in); e != nil {
          return nil, e
        }
        
        return g.ReadNum(in)
      } else if unicode.IsGraphic(c) {
        if e = unread(g, in); e != nil {
          return nil, e
        }

        return g.ReadId(in)
      }
    }

    break
  }

  return nil, g.NewError(g.Pos, "Unexpected input: %v", c)
}

func (g *G) ReadExpr(in *strings.Reader) (Form, Error) {
  g.Pos.Col++
  ef := new(ExprForm).Init()

  for {
    f, e := g.Read(in, ')')

    if e != nil {
      return nil, e
    }

    if f == nil {
      break
    }

    ef.Append(f)
  }

  return ef, nil
}

func (g *G) ReadId(in *strings.Reader) (Form, Error) {
  var buf strings.Builder
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.NewError(g.Pos, "Failed readig char: %v", re)
    }
    
    if unicode.IsSpace(c) || c == '%' || c == '(' || c == ')' {
      if e := unread(g, in); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.NewError(g.Pos, "Failed writing char: %v", we)
    }

    g.Pos.Col++
  }

  return new(IdForm).Init(g.Sym(buf.String())), nil
}

func (g *G) ReadNum(in *strings.Reader) (Form, Error) {
  var buf strings.Builder
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.NewError(g.Pos, "Error reading char: %v", re)
    }
    
    if !unicode.IsDigit(c) {
      if e := unread(g, in); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.NewError(g.Pos, "Error writing char: %v", we)
    }

    g.Pos.Col++
  }

  n, e := strconv.ParseInt(buf.String(), 10, 64)

  if e != nil {
    return nil, g.NewError(g.Pos, "Invalid num: %v", buf.String()) 
  }

  var v Val
  v.Init(g.Int, Int(n))
  return new(LitForm).Init(v), nil
}
