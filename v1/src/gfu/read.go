package gfu

import (
  "io"
  //"log"
  "strconv"
  "strings"
  "unicode"
)

func (g *G) ReadChar(pos *Pos, in *strings.Reader) (rune, E) {
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

func (g *G) Unread(pos *Pos, in *strings.Reader, c rune) E {
  if e := in.UnreadRune(); e != nil {
    return g.E(*pos, "Failed unreading char")
  }

  if c == '\n' {
    pos.Row--
  } else {
    pos.Col--
  }

  return nil
}

func (g *G) Read(pos *Pos, in *strings.Reader, out []Form, end rune) ([]Form, E) {
  var c rune
  var e E

  for {
    c, e = g.ReadChar(pos, in)

    if e != nil || c == 0 || c == end {
      return nil, e
    }
    
    switch c {
    case ' ', '\n':
      break
    case '(':
      return g.ReadExpr(pos, in, out)
    case '\'':
      return g.ReadQuote(pos, in, out, end)
    case '%':
      return g.ReadUnquote(pos, in, out, end)
    default:
      if unicode.IsDigit(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }
        
        return g.ReadNum(pos, in, out, false)
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
          return g.ReadNum(pos, in, out, true)
        }

        return g.ReadId(pos, in, out, "-")        
      } else if unicode.IsGraphic(c) {
        if e = g.Unread(pos, in, c); e != nil {
          return nil, e
        }

        return g.ReadId(pos, in, out, "")
      }

      return nil, g.E(*pos, "Unexpected input: %v", c)
    }
  }
}

func (g *G) ReadExpr(pos *Pos, in *strings.Reader, out []Form) ([]Form, E) {
  ef := new(ExprForm).Init(*pos)

  for {
    fs, e := g.Read(pos, in, ef.body, ')')

    if e != nil {
      return nil, e
    }

    if fs == nil {
      break
    }

    ef.body = fs
  }

  return append(out, ef), nil
}

func (g *G) ReadId(pos *Pos, in *strings.Reader, out []Form, prefix string) ([]Form, E) {
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

  if strings.HasSuffix(s, "..") {
    if s == ".." {
      out_len := len(out)

      if out_len == 0 {
        return nil, g.E(*pos, "Nothing to splat")
      }

      return append(out[:out_len-1], new(SplatForm).Init(fpos, out[out_len-1])), nil
    } else {
      f := new(IdForm).Init(fpos, g.Sym(s[:len(s)-2]))
      return append(out, new(SplatForm).Init(fpos, f)), nil
    }
  }
  
  return append(out, new(IdForm).Init(fpos, g.Sym(s))), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, out []Form, is_neg bool) ([]Form, E) {
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
      return nil, g.E(*pos, "Failed writing char: %v", we)
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

  if is_neg {
    n = -n
  }
  
  var v Val
  v.Init(g.IntType, int(n))
  f := new(LitForm).Init(fpos, v)
  
  if splat {
    out = append(out, new(SplatForm).Init(fpos, f))
  } else {
    out = append(out, f)
  }

  return out, nil
}

func (g *G) ReadQuote(pos *Pos, in *strings.Reader, out []Form, end rune) ([]Form, E) {
  fpos := *pos
  var fs []Form
  fs, e := g.Read(pos, in, fs, end)

  if e != nil {
    return nil, e
  }

  if len(fs) == 0 {
    return nil, g.E(*pos, "Nothing to quote")
  }

  for _, f := range fs {
    out = append(out, new(QuoteForm).Init(fpos, f))
  }

  return out, nil
}

func (g *G) ReadString(pos *Pos, in string) ([]Form, E) {
  var out []Form
  return g.Read(pos, strings.NewReader(in), out, 0)
}

func (g *G) ReadUnquote(pos *Pos, in *strings.Reader, out []Form, end rune) ([]Form, E) {
  var fs []Form
  fs, e := g.Read(pos, in, fs, end)

  if e != nil {
    return nil, e
  }

  if len(fs) == 0 {
    return nil, g.E(*pos, "Nothing to unquote")
  }

  for _, f := range fs {
    out = append(out, new(UnquoteForm).Init(f.Pos(), f))
  }

  return out, nil
}
