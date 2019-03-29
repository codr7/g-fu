package gfu

type Pos struct {
  src string
  row, col int
}

var MIN_POS Pos

func init() {
  MIN_POS.Init("n/a", 1, 0)
}

func (p *Pos) Init(src string, row, col int) *Pos {
  p.src = src
  p.row = row
  p.col = col
  return p
}
