package main

import (
  "fmt"
  "log"
  "strings"

  "./gfu"
)

func main() {
  fmt.Printf("g-fu v1.1\n\n")
  g, e := gfu.NewG()

  if e != nil {
    log.Fatal(e)
  }
  
  g.RootEnv.InitAbc(g)
  g.Debug = true

  var f gfu.Form
  pos := gfu.MIN_POS
  //f, e = g.Read(strings.NewReader("(let (x 35) ((fun (y) y) 42))"), &pos, 0)
  //f, e = g.Read(strings.NewReader("((fun (x) (+ x 7)) 35)"), &pos, 0)
  //f, e = g.Read(strings.NewReader("(let (x 35) ((fun (y) (+ x y)) 7))"), &pos, 0)
  //f, e = g.Read(strings.NewReader("(_)"), &pos, 0)
  f, e = g.Read(strings.NewReader("(bool 42)"), &pos, 0)
  //f, e = g.Read(strings.NewReader("(42 7)"), &pos, 0)

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
