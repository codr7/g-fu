package gfu

type Tag uint64

func (x Tag) Cmp(y Tag) int {
  if x < y {
    return -1
  }

  if x > y {
    return 1
  }

  return 0
}
