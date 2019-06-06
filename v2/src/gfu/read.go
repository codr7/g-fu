package gfu

import (
	"bufio"
	"fmt"
	"io"
	//"log"
	"strconv"
	"strings"
	"unicode"
)

type ERead struct {
	EBasic
	pos Pos
}

type EReadType struct {
	BasicType
}

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
		return 0, g.ERead(*pos, "Failed reading char: %v", e)
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
		return g.ERead(*pos, "Failed unreading char")
	}

	if c == '\n' {
		pos.Row--
	} else {
		pos.Col--
	}

	return nil
}

func (g *G) ReadAll(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	for {
		vs, e := g.Read(pos, in, out, "")

		if e != nil {
			return nil, e
		}

		if vs == nil {
			break
		}

		out = vs
	}

	return out, nil
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
			return nil, g.ERead(*pos, "Unexpected input: )")
		case '\'':
			return g.ReadQuote(pos, in, out, end)
		case '.':
			{
				var nc rune
				nc, e = g.ReadChar(pos, in)

				if e != nil {
					return nil, e
				}

				if nc == '.' {
					return g.ReadSplat(pos, in, out)
				}

				if !unicode.IsDigit(nc) {
					return nil, g.ERead(*pos, "Invalid input: %v", c)
				}

				if e = g.Unread(pos, in, nc); e != nil {
					return nil, e
				}

				return g.ReadNum(pos, in, out, '.')
			}
		case '%':
			return g.ReadSplice(pos, in, out, end)
		case '"':
			return g.ReadStr(pos, in, out)
		default:
			if unicode.IsDigit(c) {
				if c == '0' {
					var nc rune
					nc, e = g.ReadChar(pos, in)

					if e != nil {
						return nil, e
					}

					if nc == 'x' {
						return g.ReadByte(pos, in, out)
					}

					if e = g.Unread(pos, in, nc); e != nil {
						return nil, e
					}
				}

				return g.ReadNum(pos, in, out, c)
			} else if c == '-' {
				var nc rune
				nc, e = g.ReadChar(pos, in)

				if e != nil {
					return nil, e
				}

				is_num := unicode.IsDigit(nc) || nc == '.'

				if e = g.Unread(pos, in, nc); e != nil {
					return nil, e
				}

				if is_num {
					return g.ReadNum(pos, in, out, c)
				}

				return g.ReadId(pos, in, out, c)
			} else if unicode.IsGraphic(c) {
				if e = g.Unread(pos, in, c); e != nil {
					return nil, e
				}

				return g.ReadId(pos, in, out, 0)
			}

			return nil, g.ERead(*pos, "Invalid input: %v", c)
		}
	}
}

func (g *G) ReadByte(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	v := make([]byte, 4)
	v[0] = '0'
	v[1] = 'x'

	n, e := in.Read(v[2:])

	if e != nil {
		return nil, g.ERead(*pos, "Failed reading byte: %v", e)
	}

	sv := string(v)

	if n != 2 {
		return nil, g.E("Invalid byte: %v", sv)
	}

	b, e := strconv.ParseUint(sv, 0, 8)

	if e != nil {
		return nil, g.ERead(*pos, "Failed parsing byte: %v", e)
	}

	pos.Col += 2
	return append(out, Byte(b)), nil
}

func (g *G) ReadId(pos *Pos, in *strings.Reader, out Vec, prefix rune) (Vec, E) {
	var buf strings.Builder

	if prefix != 0 {
		buf.WriteRune(prefix)
	}

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
			return nil, g.ERead(*pos, "Failed writing char: %v", we)
		}
	}

	s := g.Sym(buf.String())

	if v := g.FindConst(s); v != nil {
		var e E

		if v, e = g.Clone(v); e != nil {
			return nil, e
		}

		return append(out, v), nil
	}

	return append(out, s), nil
}

func (g *G) ReadNum(pos *Pos, in *strings.Reader, out Vec, prefix rune) (Vec, E) {
	var buf strings.Builder

	if prefix != 0 {
		buf.WriteRune(prefix)
	}

	is_float := prefix == '.'

	for {
		c, e := g.ReadChar(pos, in)

		if e != nil {
			return nil, e
		}

		if c == 0 {
			break
		}

		is_float = is_float || c == '.'

		if !unicode.IsDigit(c) && c != '.' {
			if e := g.Unread(pos, in, c); e != nil {
				return nil, e
			}

			break
		}

		if _, we := buf.WriteRune(c); we != nil {
			return nil, g.ERead(*pos, "Failed writing char: %v", we)
		}
	}

	s := buf.String()
	rs := []rune(s)

	if rs[0] == '-' && len(rs) > 1 && rs[1] == '.' {
		s = fmt.Sprintf("-0%v", string(rs[1:]))
	}

	if is_float {
		var v Float
		e := v.Parse(g, s)

		if e != nil {
			return nil, e
		}

		return append(out, v), nil
	}

	n, e := strconv.ParseInt(s, 10, 64)

	if e != nil {
		return nil, g.ERead(*pos, "Invalid Int: %v", s)
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
		return nil, g.ERead(vpos, "Nothing to quote")
	}

	return append(out, NewQuote(g, vs[0])), nil
}

func (g *G) ReadSplat(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	i := len(out)

	if i == 0 {
		return nil, g.ERead(*pos, "Missing splat value")
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
		return nil, g.ERead(vpos, "Nothing to eval")
	}

	return append(out, NewSplice(g, vs[0])), nil
}

func (g *G) ReadStr(pos *Pos, in *strings.Reader, out Vec) (Vec, E) {
	var buf strings.Builder
	pc := ' '

	for {
		c, e := g.ReadChar(pos, in)

		if e != nil {
			return nil, e
		}

		if c == 0 || c == '"' {
			break
		}

		if c == '\\' {
			if pc == '\\' {
				pc = ' '
				continue
			}

			nc, e := g.ReadChar(pos, in)

			if e != nil {
				return nil, e
			}

			switch nc {
			case '"':
				c = nc
			case 'e':
				c = '\x1b'
			case 'n':
				c = '\n'
			}
		}

		if _, we := buf.WriteRune(c); we != nil {
			return nil, g.ERead(*pos, "Failed writing char: %v", we)
		}

		pc = c
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

func (e *ERead) Init(g *G, pos Pos, msg string) *ERead {
	e.EBasic.Init(g, msg)
	e.pos = pos
	return e
}

func (_ ERead) Type(g *G) Type {
	return &g.EReadType
}

func (_ EReadType) Dump(g *G, val Val, out *bufio.Writer) E {
	e := val.(*ERead)
	p := &e.pos

	fmt.Fprintf(out,
		"Read error in '%s'; row %v, col %v:\n%v",
		p.src, p.Row, p.Col, e.msg)

	return nil
}

func (g *G) ERead(pos Pos, msg string, args ...interface{}) *ERead {
	msg = fmt.Sprintf(msg, args...)
	e := new(ERead).Init(g, pos, msg)
	return e
}
