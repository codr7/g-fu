package main

import (
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"
  
  "./gfu"
)

var prof = flag.String("prof", "", "Write CPU profile to specified file")

func main() {
  fmt.Printf("g-fu v1.2\n\n")
  g, e := gfu.NewG()

  if e != nil {
    log.Fatal(e)
  }
  
  g.RootEnv.InitAbc(g)
  g.Debug = true
  pos := gfu.MIN_POS
  flag.Parse()
  
  if *prof != "" {
    f, e := os.Create(*prof)

    if e != nil {
      log.Fatal(e)
    }
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()    
  }

  for _, a := range flag.Args() {
    if _, e := g.Load(a, &g.RootEnv, pos); e != nil {
      log.Fatal(e);
    }
  }  
}
