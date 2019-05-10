package gfu

import (
	"strings"
)

type Dumper interface {
	Dump(*strings.Builder)
}

func DumpString(d Dumper) string {
	var out strings.Builder
	d.Dump(&out)
	return out.String()
}
