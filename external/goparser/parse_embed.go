// This is literally a copy of go/parser.ParseDir with minimal changes to 
// allow it to support EmbedFS directly. I imagine this wouldn't be accepted
// back into base go as it depends on embed which was only introduced in 1.16.
// Also no tests as yet.
//
package goparser

import (
	"embed"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

// ParseEmbedFSDir calls ParseFile for all files with names ending in ".go"
// in the directory specified by path in embedFS, and returns a map of
// package name -> package AST with all the packages found.
//
// If filter != nil, only the files with fs.FileInfo entries passing through
// the filter (and ending in ".go") are considered. The mode bits are passed
// to ParseFile unchanged. Position information is recorded in fset, which
// must not be nil.
//
// If the directory couldn't be read, a nil map and the respective error are
// returned. If a parse error occurred, a non-nil but incomplete map and the
// first error encountered are returned.
//
func ParseEmbedFSDir(fset *token.FileSet, path string, embedFS embed.FS, filter func(fs.FileInfo) bool, mode parser.Mode) (pkgs map[string]*ast.Package, first error) {
	list, err := embedFS.ReadDir(path)
	if err != nil {
		return nil, err
	}

	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			continue
		}
		if filter != nil {
			info, err := d.Info()
			if err != nil {
				return nil, err
			}
			if !filter(info) {
				continue
			}
		}
		filename := filepath.Join(path, d.Name())
		fileBytes, err := embedFS.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		if src, err := parser.ParseFile(fset, filename, fileBytes, mode); err == nil {
			name := src.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:  name,
					Files: make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[filename] = src
		} else if first == nil {
			first = err
		}
	}

	return
}
