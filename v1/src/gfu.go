package main

import (
  "bufio"
  "flag"
  "fmt"
  "log"
  "os"
  "runtime/pprof"
  "strings"
  
  "./gfu"
)

var prof = flag.String("prof", "", "Write CPU profile to specified file")

func main() {
  g, e := gfu.NewG()

  if e != nil {
    log.Fatal(e)
  }
  
  g.RootEnv.InitAbc(g)
  g.Debug = true
  flag.Parse()
  
  if *prof != "" {
    f, e := os.Create(*prof)

    if e != nil {
      log.Fatal(e)
    }
    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()    
  }

  args := flag.Args()
  
  if len(args) == 0 {
    fmt.Printf("g-fu v1.3\n\nPress Return twice to evaluate.\n\n  ")
    in := bufio.NewScanner(os.Stdin)
    var buf strings.Builder
    
    for in.Scan() {
      line := in.Text()

      if len(line) == 0 {
        v, e := g.EvalString(gfu.MIN_POS, buf.String(), &g.RootEnv)
        buf.Reset()

        if e == nil {
          fmt.Printf("\r%v\n", v)
        } else {
          fmt.Printf("\r%v\n", e)
        }
      } else {
        buf.WriteString(line)
      }

      fmt.Printf("  ")
    }

    if e := in.Err(); e != nil {
      log.Fatal(e)
    }
  } else {
    for _, a := range args {
      if _, e := g.Load(gfu.MIN_POS, a, &g.RootEnv); e != nil {
        log.Fatal(e);
      }
    }
  }
}
