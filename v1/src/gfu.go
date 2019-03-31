package main

import (
  "fmt"
  "log"
  "os"
  
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
  pos := gfu.MIN_POS

  if len(os.Args) > 1 {
    args := os.Args[1:]

    if _, e := g.Load(args[0], &g.RootEnv, pos); e != nil {
      log.Fatal(e);
    }
  }  
}
