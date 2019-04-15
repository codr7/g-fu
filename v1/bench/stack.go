package main

import (
  "fmt"
  "./util"
)

func main() {
  fmt.Printf("%v\n", util.Bench(10, func() {
    var s []int64
    
    for i := int64(0); i < 100000; i++ {
      s = append(s, i)
    }

    for i := int64(0); i < 100000; i++ {
      s = s[:len(s)-1]
    }
  }))
}
