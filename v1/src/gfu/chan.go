package gfu

type Chan chan Val

func NewChan(max_buf Int) Chan {
  return make(Chan, max_buf)
}
