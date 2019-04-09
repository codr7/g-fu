package gfu

type Pos struct {
  src string
  Row, Col int
}

var (
  INIT_POS, NIL_POS Pos
)

func init() {
  INIT_POS.Init("n/a")
  
  NIL_POS.src = "n/a"
  NIL_POS.Row = -1
  NIL_POS.Col = -1
}

func (p *Pos) Init(src string) *Pos {
  p.src = src
  p.Row = 1
  p.Col = 0
  return p
}
