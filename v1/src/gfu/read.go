package gfu

import (
	"io"
	//"log"
	"strconv"
	"strings"
	"unicode"
)

type CharSet string

func (s CharSet) Member(c rune) bool {
	return strings.IndexRune(string(s), c) != -1
}

func (g *G) ReadChar(pos *Pos, in *strings.Reader) (rune, E) {
	c, _, e := in.ReadRune()

	if e == io.EOF {
		return 0, nil
	}

	if e != nil {
		return 0, g.ReadE(*pos, "Failed reading char: %v", e)
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
		return g.ReadE(*pos, "Failed unreading char")
	}

	if c == '\n' {
		pos.Row--
	} else {
		pos.Col--
	}

	return nil
}

func (g *G) Read(pos *Pos, in *strings.Reader, out Vec, end CharSet) (Vec, E) {
	var c rune
	var e E

	for {
		c, e = g.ReadChar(pos, in)

		if e != nil {
			return nil, e
		}

		if end.Member(c) {
			if e = g.Unread(pos, in, c); e != nil {
				return nil, e
			}

			c = 0
		}

		if c == 0 {
			return nil, e
		}

		switch c {
		case ' ', '\n':
			break
		case '(':
			return g.ReadVec(pos, in, out)
		case ')':
			return nil, g.ReadE(*pos, "Unexpected input: )")
		case '\'':
			return g.ReadQuote(pos, in, out, end)
		case '.':
			return g.ReadSplat(pos, in, out)
		case '%':
			return g.ReadSplice(pos, in, out, end)
		case '"':
			return g.ReadStr(pos, in, out)
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

				is_num := unicode.IsDigit(nc)

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

			return nil, g.ReadE(*pos, "Unexpected input: %v", c)
		}
	}
}

func (g *G) ReadId(pos *Pos, in *strings.Reader, out Vec, prefix string) (Vec, E) {
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
			c == '.' || c == '%' || c == '(' || c == ')' {
			if e := g.Unread(pos, in, c); e != nil {
				return nil, e
			}

			break
		}

		if _, we := buf.WriteRune(c); we != nil {
			return nil, g.ReadE(*pos, "Failed writing char: %v", we)
		}
	}

	s := g.Sym(buf.String())

	if v := g.FindConst(s); v != nil {
		var e E

		if v, e = v.Clone(g); e != nil {
			return nil, e
		}

		return append(out, v), nil
	}

	return append(out, s), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, out Vec, is_neg bool) (Vec, E) {
	var buf strings.Builder

	for {
		c, e := g.ReadChar(pos, in)

		if e != nil {
			return nil, e
		}

		if c == 0 {
			break
		}

		if !unicode.IsDigit(c) {
			if e := g.Unread(pos, in, c); e != nil {
				return nil, e
			}

			break
		}

		if _, we := buf.WriteRune(c); we != nil {
			return nil, g.ReadE(*pos, "Failed writing char: %v", we)
		}
	}

	s := buf.String()
	n, e := strconv.ParseInt(s, 10, 64)

	if e != nil {
		return nil, g.ReadE(*pos, "Invalid num: %v", s)
	}

	if is_neg {
		n = -n
	}

	return append(out, Int(n)), nil
}

func (g *G) ReadQuote(pos *Pos, in *strings.Reader, out Vec, end CharSet) (Vec, E) {
	vpos := *pos
	vs, e := g.Read(pos, in, nil, end)

	if e != nil {
		return nil, e
	}

	if len(vs) == 0 {
		return nil, g.ReadE(vpos, "Nothing to quote")
	}

	return append(out, NewQuote(g, vs[0])), nil
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
		return nil, g.ReadE(*pos, "Invalid input: .%v", string(nc))
	}

	i := len(out)

	if i == 0 {
		return nil, g.ReadE(*pos, "Missing splat value")
	}

	v := out[i-1]
	return append(out[:i-1], NewSplat(g, v)), nil
}

func (g *G) ReadSplice(pos *Pos, in *strings.Reader, out Vec, end CharSet) (Vec, E) {
	vpos := *pos
	vpos.Col--

	vs, e := g.Read(pos, in, nil, end)

	if e != nil {
		return nil, e
	}

	if len(vs) == 0 {
		return nil, g.ReadE(vpos, "Nothing to eval")
	}

	return append(out, NewSplice(g, vs[0])), nil
}

func (g *G) ReadStr(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	var buf strings.Builder

	for {
		c, e := g.ReadChar(pos, in)

		if e != nil {
			return nil, e
		}

		if c == 0 || c == '"' {
			break
		}

		if _, we := buf.WriteRune(c); we != nil {
			return nil, g.ReadE(*pos, "Failed writing char: %v", we)
		}
	}

	return append(out, Str(buf.String())), nil
}

func (g *G) ReadVec(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	var body Vec

	for {
		vs, e := g.Read(pos, in, body, ")")

		if e != nil {
			return nil, e
		}

		if vs == nil {
			break
		}

		body = vs
	}

	c, e := g.ReadChar(pos, in)

	if e != nil {
		return nil, e
	}

	if c != ')' {
		return nil, g.E("Invalid vec end: %v", string(c))
	}

	return append(out, body), nil
}
