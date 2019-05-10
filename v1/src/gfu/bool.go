package gfu

func (g *G) BoolVal(val Val) (Val, E) {
	bv, e := g.Bool(val)

	if e != nil {
		return nil, e
	}

	if bv {
		return &g.T, nil
	}

	return &g.F, nil
}
