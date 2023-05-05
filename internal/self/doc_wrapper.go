package self

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"

	"github.com/azekeil/grec/external/goparser"
	"github.com/azekeil/grec/internal"
)

// parseEmbeds parses source files from an embed
func parseEmbeds(path string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet() // positions are relative to fset

	d, err := goparser.ParseEmbedFSDir(fset, path, internal.CommandsEmbedFS, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func getDocPackage(d map[string]*ast.Package, name string) *doc.Package {
	for k, f := range d {
		if k == name {
			return doc.New(f, "./", 0)
		}
	}
	return nil
}

func getDocMethods(p *doc.Package, typeName string) DocFuncs {
	s := make(map[string]*doc.Func, len(p.Types))
	for _, t := range p.Types {
		if t.Name == typeName {
			for _, m := range t.Methods {
				s[m.Name] = m
			}
		}
	}
	return DocFuncs(s)
}

func MakeHelp(path, pkg, typ string) map[string]*doc.Func {
	d, err := parseEmbeds(path)
	if err != nil {
		panic(err)
	}
	return getDocMethods(getDocPackage(d, pkg), typ)
}
