package self

import (
	"go/doc"
	"strings"
)

type DocFuncs map[string]*doc.Func

func (f DocFuncs) Summary(name string) string {
	return strings.Split(f[name].Doc, "\n")[0]
}

func (f DocFuncs) CommandHelp(name string) string {
	return f[name].Doc
}

func (f DocFuncs) AllSummaries() []string {
	var s []string
	for k := range f {
		s = append(s, f.Summary(k))
	}
	return s
}
