package main

import (
  "fmt"
  "log"

  "./gfu"
)

func main() {
  fmt.Printf("g-fu v1.2\n\n")
  g, e := gfu.NewG()

  if e != nil {
    log.Fatal(e)
  }
  
  g.RootEnv.InitAbc(g)
  g.Debug = true

  var f gfu.Form
  pos := gfu.MIN_POS

  f, e = g.ReadString("(let (fib (fun (n) (or (and n (+ (fib (- n 1)) 1)) 1))) (fib 7))", &pos)
  //f, e = g.ReadString("(and T 42)", &pos)
  //f, e = g.ReadString("(or 42 T _)", &pos)
  //f, e = g.ReadString("(do 1 2 3)", &pos)
  //f, e = g.ReadString("(+ 42 _..)", &pos)
  //f, e = g.ReadString("(+ 7..)", &pos)
  //f, e = g.ReadString("(- 42 1 2 3)", &pos)
  //f, e = g.ReadString("((fun (xs..) xs) 1 2 3)", &pos)
  //f, e = g.ReadString("((fun (xs..) (+ xs..)) 1 2 3)", &pos)
  //f, e = g.ReadString("(let (x 35) ((fun (y) y) 42))", &pos)
  //f, e = g.ReadString("((fun (x) (+ x 7)) 35)", &pos)
  //f, e = g.ReadString("(let (x 35) ((fun (y) (+ x y)) 7))", &pos)
  //f, e = g.ReadString("(_)", &pos)
  //f, e = g.ReadString("(bool 42)", &pos)
  //f, e = g.ReadString("(42 7)", &pos)

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
