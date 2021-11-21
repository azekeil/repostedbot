package self

import (
	"go/doc"
	"sort"
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
	sort.Strings(s)
	return s
}

// Capitalise does a case-insensitive comparison on the function names
// and returns the correctly-capitalised name if present
func (f DocFuncs) Capitalise(name string) string {
	lname := strings.ToLower(name)
	for _, k := range f {
		if lname == strings.ToLower(k.Name) {
			return k.Name
		}
	}
	return ""
}
