package self

import (
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
)

func ParseDir(path string) (map[string]*ast.Package, error) {
	fset := token.NewFileSet() // positions are relative to fset

	d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return d, nil
}

func GetDocPackage(d map[string]*ast.Package, name string) *doc.Package {
	for k, f := range d {
		if k == name {
			return doc.New(f, "./", 0)
		}
	}
	return nil
}

func GetDocMethods(p *doc.Package, typeName string) DocFuncs {
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
