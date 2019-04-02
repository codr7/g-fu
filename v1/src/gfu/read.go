package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func (g *G) Unread(pos Pos, in *strings.Reader) Error {
  if e := in.UnreadRune(); e != nil {
    return g.E(pos, "Error unreading char")
  }

  return nil
}

func (g *G) Read(pos *Pos, in *strings.Reader, end rune) (Form, Error) {
  var c rune
  var e Error

  for {
    var re error
    c, _, re = in.ReadRune()

    if re == io.EOF {
      return nil, nil
    }
    
    if re != nil {
      return nil, g.E(*pos, "Failed reading char: %v", re)
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
      return g.ReadExpr(pos, in)
    default:
      if unicode.IsDigit(c) {
        if e = g.Unread(*pos, in); e != nil {
          return nil, e
        }
        
        return g.ReadNum(pos, in, false)
      } else if c == '-' {
        pos.Col++
        nc, _, re := in.ReadRune()
        
        if re != nil {
          return nil, g.E(*pos, "Failed reading char: %v", re)
        }

        is_num := unicode.IsDigit(nc);
        
        if e = g.Unread(*pos, in); e != nil {
          return nil, e
        }
          
        if is_num {
          return g.ReadNum(pos, in, true)
        }

        return g.ReadId(pos, in, "-")        
      } else if unicode.IsGraphic(c) {
        if e = g.Unread(*pos, in); e != nil {
          return nil, e
        }

        return g.ReadId(pos, in, "")
      }
    }

    break
  }

  return nil, g.E(*pos, "Unexpected input: %v", c)
}

func (g *G) ReadExpr(pos *Pos, in *strings.Reader) (Form, Error) {
  ef := new(ExprForm).Init(*pos)
  pos.Col++

  for {
    f, e := g.Read(pos, in, ')')

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

func (g *G) ReadId(pos *Pos, in *strings.Reader, prefix string) (Form, Error) {
  fpos := *pos
  var buf strings.Builder
  buf.WriteString(prefix)
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.E(*pos, "Failed readig char: %v", re)
    }
    
    if unicode.IsSpace(c) || c == '%' || c == '(' || c == ')' {
      if e := g.Unread(*pos, in); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.E(*pos, "Failed writing char: %v", we)
    }

    pos.Col++
  }

  s := buf.String()

  if s == ".." {
    return new(SplatForm).Init(fpos), nil
  }
  
  return new(IdForm).Init(fpos, g.S(buf.String())), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, is_neg bool) (Form, Error) {
  fpos := *pos
  var buf strings.Builder
  
  for {
    c, _, re := in.ReadRune()

    if re != nil {
      if re == io.EOF {
        break
      }
      
      return nil, g.E(*pos, "Error reading char: %v", re)
    }
    
    if !unicode.IsDigit(c) && c != '.' {
      if e := g.Unread(*pos, in); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.E(*pos, "Error writing char: %v", we)
    }

    pos.Col++
  }

  s := buf.String()
  splat := false
  
  if strings.HasSuffix(s, "..") {
    s = s[:len(s)-2]
    splat = true
  }
  
  n, e := strconv.ParseInt(s, 10, 64)

  if e != nil {
    return nil, g.E(*pos, "Invalid num: %v", s) 
  }

  var v Val

  if is_neg {
    n = -n
  }
  
  v.Init(g.Int, Int(n))

  if splat {
    v.Init(g.Splat, v)
  }
  
  return new(LitForm).Init(fpos, v), nil
}

func (g *G) ReadString(pos *Pos, in string) (Form, Error) {
  return g.Read(pos, strings.NewReader(in), 0)
}
