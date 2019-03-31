package gfu

type Pos struct {
  src string
  Row, Col int
}

var MIN_POS Pos

func init() {
  MIN_POS.Init("n/a")
}

func (p *Pos) Init(src string) *Pos {
  p.src = src
  p.Row = 1
  p.Col = 0
  return p
}
