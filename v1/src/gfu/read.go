package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func (g *G) Unread(in *strings.Reader, pos Pos) Error {
  if e := in.UnreadRune(); e != nil {
    return g.NewError(pos, "Error unreading char")
  }

  return nil
}

func (g *G) Read(in *strings.Reader, pos *Pos, end rune) (Form, Error) {
  var c rune
  var e Error

  for {
    var re error
    c, _, re = in.ReadRune()

    if re == io.EOF {
      return nil, nil
    }
    
    if re != nil {
      return nil, g.NewError(*pos, "Error reading char: %v", re)
    }
    
    if c == end {
      return nil, nil
    }
    
    switch c {
    case ' ':
      pos.Col++
      continue
    case '\n':
      pos.Col = MIN_POS.Col
      pos.Row++
      continue
    case '(':
      return g.ReadExpr(in, pos)
    default:
      if unicode.IsDigit(c) {
        if e = g.Unread(in, *pos); e != nil {
          return nil, e
        }
        
        return g.ReadNum(in, pos)
      } else if unicode.IsGraphic(c) {
        if e = g.Unread(in, *pos); e != nil {
          return nil, e
        }

        return g.ReadId(in, pos)
      }
    }

    break
  }

  return nil, g.NewError(*pos, "Unexpected input: %v", c)
}

func (g *G) ReadExpr(in *strings.Reader, pos *Pos) (Form, Error) {
  ef := new(ExprForm).Init(*pos)
  pos.Col++

  for {
    f, e := g.Read(in, pos, ')')

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

func (g *G) ReadId(in *strings.Reader, pos *Pos) (Form, Error) {
  fpos := *pos
  var buf strings.Builder
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.NewError(*pos, "Failed readig char: %v", re)
    }
    
    if unicode.IsSpace(c) || c == '%' || c == '(' || c == ')' {
      if e := g.Unread(in, *pos); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.NewError(*pos, "Failed writing char: %v", we)
    }

    pos.Col++
  }

  return new(IdForm).Init(fpos, g.Sym(buf.String())), nil
}

func (g *G) ReadNum(in *strings.Reader, pos *Pos) (Form, Error) {
  fpos := *pos
  var buf strings.Builder
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.NewError(*pos, "Error reading char: %v", re)
    }
    
    if !unicode.IsDigit(c) {
      if e := g.Unread(in, *pos); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.NewError(*pos, "Error writing char: %v", we)
    }

    pos.Col++
  }

  n, e := strconv.ParseInt(buf.String(), 10, 64)

  if e != nil {
    return nil, g.NewError(*pos, "Invalid num: %v", buf.String()) 
  }

  var v Val
  v.Init(g.Int, Int(n))
  return new(LitForm).Init(fpos, v), nil
}
