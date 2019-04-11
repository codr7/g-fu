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

func (g *G) Read(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
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
      return g.ReadVec(pos, in, out)
    case '\'':
      return g.ReadQuote(pos, in, out, end)
    case '?':
      return g.ReadOpt(pos, in, out)
    case '.':
      return g.ReadSplat(pos, in, out)
    case '%':
      return g.ReadSplice(pos, in, out, end)
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

func (g *G) ReadVec(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  vpos := *pos
  var body Vec

  for {
    vs, e := g.Read(pos, in, body, ')')

    if e != nil {
      return nil, e
    }

    if vs == nil {
      break
    }

    body = vs
  }

  var v Val
  v.Init(vpos, g.VecType, body)
  return append(out, v), nil
}

func (g *G) ReadId(pos *Pos, in *strings.Reader, out Vec, prefix string) (Vec, E) {
  vpos := *pos
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

    if unicode.IsSpace(c) ||
      c == '.' || c == '?' || c == '%' || c == '(' || c == ')' {
      if e := g.Unread(pos, in, c); e != nil {
        return nil, e
      }

      break
    }

    if _, we := buf.WriteRune(c); we != nil {
      return nil, g.E(*pos, "Failed writing char: %v", we)
    }
  }

  var v Val
  v.Init(vpos, g.SymType, g.Sym(buf.String()))
  return append(out, v), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, out Vec, is_neg bool) (Vec, E) {
  vpos := *pos
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
  v.Init(vpos, g.IntType, int(n))
  
  if splat {
    v.Init(vpos, g.SplatType, v)
  }

  return append(out, v), nil

  return out, nil
}

func (g *G) ReadOpt(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  i := len(out)
  
  if i == 0 {
    return nil, g.E(*pos, "Missing opt value")        
  }
  
  v := &out[i-1]
  v.Init(*pos, g.OptType, *v)
  return out, nil      
}

func (g *G) ReadQuote(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
  vpos := *pos
  var vs Vec
  vs, e := g.Read(pos, in, vs, end)

  if e != nil {
    return nil, e
  }

  if len(vs) == 0 {
    return nil, g.E(vpos, "Nothing to quote")
  }

  v := vs[0]
  v.Init(vpos, g.QuoteType, v)
  return append(out, v), nil
}

func (g *G) ReadSplat(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
  vpos := *pos
  vpos.Col--
  
  var nc rune
  var e E
  
  nc, e = g.ReadChar(pos, in)
  
  if e != nil {
    return nil, e
  }

  if nc != '.' {
    return nil, g.E(*pos, "Invalid input: .%v", nc)
  }
  
  i := len(out)

  if i == 0 {
    return nil, g.E(*pos, "Missing splat value")        
  }

  v := &out[i-1]
  v.Init(vpos, g.SplatType, *v)
  return out, nil      
}

func (g *G) ReadSplice(pos *Pos, in *strings.Reader, out Vec, end rune) (Vec, E) {
  vpos := *pos
  vpos.Col--
  
  var vs Vec
  vs, e := g.Read(pos, in, vs, end)

  if e != nil {
    return nil, e
  }

  if len(vs) == 0 {
    return nil, g.E(vpos, "Nothing to eval")
  }

  v := vs[0]
  v.Init(vpos, g.SpliceType, v)
  return append(out, v), nil
}

