package gfu

import (
  "bufio"
  "strings"
)

type Dumper interface {
  Dump(*G, *bufio.Writer) E
}

func (g *G) DumpString(d Dumper) (string, E) {
  var out strings.Builder
  w := bufio.NewWriter(&out)
  
  if e := d.Dump(g, w); e != nil {
    return "", e
  }

  w.Flush()
  return out.String(), nil
}
