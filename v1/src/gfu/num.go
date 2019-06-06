package gfu

type Num interface {
	Type
	Abs(*G, Val) (Val, E)
	Add(*G, Val, Val) (Val, E)
	Byte(*G, Val) (Byte, E)
	Div(*G, Val, Val) (Val, E)
	Float(*G, Val) (Float, E)
	Int(*G, Val) (Int, E)
	Mul(*G, Val, Val) (Val, E)
	Neg(*G, Val) (Val, E)
	Sub(*G, Val, Val) (Val, E)
}

type NumType struct {
	BasicType
}

func (_ *NumType) Abs(g *G, val Val) (Val, E) {
	return nil, g.E("Abs not supported: %v", val.Type(g))
}

func (_ *NumType) Add(g *G, x Val, _ Val) (Val, E) {
	return nil, g.E("Add not supported: %v", x.Type(g))
}

func (_ *NumType) Byte(g *G, val Val) (Byte, E) {
	return Byte(0), g.E("Byte not supported: %v", val.Type(g))
}

func (_ *NumType) Div(g *G, x Val, _ Val) (Val, E) {
	return nil, g.E("Div not supported: %v", x.Type(g))
}

func (_ *NumType) Float(g *G, val Val) (Float, E) {
	return Float{}, g.E("Float not supported: %v", val.Type(g))
}

func (_ *NumType) Inc(g *G, val, delta Val) (Val, E) {
	return nil, g.E("Inc not supported: %v", val.Type(g))
}

func (_ *NumType) Int(g *G, val Val) (Int, E) {
	return 0, g.E("Int not supported: %v", val.Type(g))
}

func (_ *NumType) Mul(g *G, x Val, _ Val) (Val, E) {
	return nil, g.E("Mul not supported: %v", x.Type(g))
}

func (_ *NumType) Neg(g *G, val Val) (Val, E) {
	return nil, g.E("Neg not supported: %v", val.Type(g))
}

func (_ *NumType) Sub(g *G, x Val, _ Val) (Val, E) {
	return nil, g.E("Sub not supported: %v", x.Type(g))
}

func (g *G) Abs(x Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Abs(g, x)
}

func (g *G) Add(x, y Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Add(g, x, y)
}

func (g *G) Byte(val Val) (Byte, E) {
	t := val.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return Byte(0), g.E("Expected Num: %v", t)
	}

	return nt.Byte(g, val)
}

func (g *G) Div(x, y Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Div(g, x, y)
}

func (g *G) Float(x Val) (Float, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return Float{}, g.E("Expected Num: %v", t)
	}

	return nt.Float(g, x)
}

func (g *G) Int(x Val) (Int, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return Int(0), g.E("Expected Num: %v", t)
	}

	return nt.Int(g, x)
}

func (g *G) Mul(x, y Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Mul(g, x, y)
}

func (g *G) Neg(x Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Neg(g, x)
}

func (g *G) Sub(x, y Val) (Val, E) {
	t := x.Type(g)
	nt, ok := t.(Num)

	if !ok {
		return nil, g.E("Expected Num: %v", t)
	}

	return nt.Sub(g, x, y)
}
