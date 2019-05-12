package gfu

import (
  "strings"
)

type Dumper interface {
  Dump(*G, *strings.Builder) E
}

func DumpString(g *G, d Dumper) (string, E) {
  var out strings.Builder

  if e := d.Dump(g, &out); e != nil {
    return "", e
  }
  
  return out.String(), nil
}
