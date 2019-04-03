package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func (g *G) ReadChar(pos *Pos, in *strings.Reader) (rune, Error) {
  c, _, e := in.ReadRune()

  if e == io.EOF {
    return 0, nil
  }
    
  if e != nil {
    return 0, g.E(*pos, "Failed reading char: %v", e)
  }

  if c == '\n' {
    pos.Col = INIT_POS.Col
    pos.Row++
  } else {
    pos.Col++
  }
  
  return c, nil
}

func (g *G) Unread(pos *Pos, in *strings.Reader, c rune) Error {
  if e := in.UnreadRune(); e != nil {
    return g.E(*pos, "Error unreading char")
  }

  if c == '\n' {
    pos.Row--
  } else {
    pos.Col--
  }

  return nil
}

func (g *G) Read(pos *Pos, in *strings.Reader, end rune) (Form, Error) {
  var c rune
  var e Error

  for {
    c, e = g.ReadChar(pos, in)
    
    if e != nil || c == 0 || c == end {
      return nil, e
    }
    
    switch c {
    case ' ', '\n':
      break
    case '(':
      return g.ReadExpr(pos, in)
    case '\'':
      return g.ReadQuote(pos, in, end)
    default:
      if unicode.IsDigit(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }
        
        return g.ReadNum(pos, in, false)
      } else if c == '-' {
        var nc rune
        nc, e = g.ReadChar(pos, in)
        
        if e != nil {
          return nil, e
        }

        is_num := unicode.IsDigit(nc);
        
        if e = g.Unread(pos, in, nc); e != nil {
          return nil, e
        }
          
        if is_num {
          return g.ReadNum(pos, in, true)
        }

        return g.ReadId(pos, in, "-")        
      } else if unicode.IsGraphic(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }

        return g.ReadId(pos, in, "")
      }

      return nil, g.E(*pos, "Unexpected input: %v", c)
    }
  }
}

func (g *G) ReadExpr(pos *Pos, in *strings.Reader) (Form, Error) {
  ef := new(ExprForm).Init(*pos)

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
    c, e := g.ReadChar(pos, in)

    if e != nil {      
      return nil, e
    }

    if c == 0 {
      break
    }

    if unicode.IsSpace(c) || c == '%' || c == '(' || c == ')' {
      if e := g.Unread(pos, in, c); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.E(*pos, "Failed writing char: %v", we)
    }
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
    c, e := g.ReadChar(pos, in)

    if e != nil {      
      return nil, e
    }

    if c == 0 {
      break
    }
    
    if !unicode.IsDigit(c) && c != '.' {
      if e := g.Unread(pos, in, c); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.E(*pos, "Error writing char: %v", we)
    }
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

func (g *G) ReadQuote(pos *Pos, in *strings.Reader, end rune) (Form, Error) {
  fpos := *pos
  f, e := g.Read(pos, in, end)

  if e != nil {
    return nil, e
  }

  return new(QuoteForm).Init(fpos, f), nil
}

func (g *G) ReadString(pos *Pos, in string) (Form, Error) {
  return g.Read(pos, strings.NewReader(in), 0)
}
