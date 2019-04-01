package main

import (
  "fmt"
  "./util"
)

func fib(n, a, b int) int {
  switch n {
  case 0:
    return a
  case 1:
    return b
  default:
    break
  }

  return fib(n-1, b, a+b)
}

func main() {
  fmt.Printf("%v\n", util.Bench(10, func() {
    for i := 0; i < 10000; i++ {
      fib(20, 0, 1)
    }
  }))
}
