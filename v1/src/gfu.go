package main

import (
  "fmt"
  "log"
  "strings"

  "./gfu"
)

func main() {
  fmt.Printf("g-fu v1\n\n")
  g, e := gfu.NewG()

  if e != nil {
    log.Fatal(e)
  }
  
  g.Debug = true

  g.Pos = gfu.MIN_POS
  var f gfu.Form
  f, e = g.Read(strings.NewReader("((fun () 42))"), 0)

  if e != nil {
    log.Fatal(e)
  }

  if f == nil {
    log.Fatalf("Missing form")
  }
  
  fmt.Printf("%v\n", f)
  var result gfu.Val
  
  result, e = f.Eval(g, &g.RootEnv)

  if e != nil {
    log.Fatal(e)
  }

  fmt.Printf("%v\n", result)
}
