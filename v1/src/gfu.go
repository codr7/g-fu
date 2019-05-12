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
  g.RootEnv.InitMath(g)
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
  env := g.NewEnv()

  if len(args) == 0 {
    fmt.Printf("g-fu v1.13\n\nPress Return twice to evaluate.\n\n  ")
    in := bufio.NewScanner(os.Stdin)
    var buf strings.Builder

    for in.Scan() {
      line := in.Text()

      if len(line) == 0 {
        if buf.Len() > 0 {
          v, e := g.EvalString(&g.MainTask, env, gfu.INIT_POS, buf.String())

          if e == nil {
            if s, e := g.String(v); e != nil {
              log.Fatal(e)
            } else {
              fmt.Printf("\r%v\n", s)
            }
          } else {
            fmt.Printf("\r%v\n", e)
          }
        }

        buf.Reset()
      } else {
        buf.WriteString(line)
        buf.WriteRune('\n')
      }

      fmt.Printf("  ")
    }

    if e := in.Err(); e != nil {
      log.Fatal(e)
    }
  } else {
    for _, a := range args {
      if _, e := g.Load(&g.MainTask, env, a); e != nil {
        log.Fatal(e)
      }
    }
  }
}
