package main

import (
  "fmt"
  "./util"
)

func fib(n int) int {
  if n < 2 {
    return n
  }

  return fib(n-1) + fib(n-2)
}

func main() {
  fmt.Printf("%v\n", util.Bench(10, func() {
    for i := 0; i < 10; i++ {
      fib(20)
    }
  }))
}
