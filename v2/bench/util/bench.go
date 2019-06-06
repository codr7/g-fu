package util

import (
  "time"
)

func Bench(n int, body func()) int64 {
  for i := 0; i < n; i++ {
    body()
  }

  t := time.Now()
  
  for i := 0; i < n; i++ {
    body()
  }

  return time.Now().Sub(t).Nanoseconds() / 1000000
}
