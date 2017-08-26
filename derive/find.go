//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package derive

import (
	"go/ast"
	"go/types"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

const derivedFilename = "derived.gen.go"

type fileInfo struct {
	astFile   *ast.File
	fullpath  string
	undefined []*call
	derived   []*call
	funcNames map[string]struct{}
}

func newFileInfos(program *loader.Program, pkgInfo *loader.PackageInfo) []*fileInfo {
	files := []*fileInfo{}
	for i := range pkgInfo.Files {
		astFile := pkgInfo.Files[i]
		file := program.Fset.File(astFile.Pos())
		if file == nil {
			// probably derived.gen.go has non parsable code.
			continue
		}
		fullpath := file.Name()

		_, fname := filepath.Split(fullpath)
		if fname == derivedFilename {
			continue
		}

		f := &finder{program, pkgInfo, nil, nil, make(map[string]struct{})}
		for _, d := range astFile.Decls {
			ast.Walk(f, d)
		}
		undefined := make([]*call, len(f.undefined))
		for i := range f.undefined {
			undefined[i] = newCall(pkgInfo, f.undefined[i])
		}
		derived := make([]*call, len(f.derived))
		for i := range f.derived {
			derived[i] = newCall(pkgInfo, f.derived[i])
		}

		files = append(files, &fileInfo{
			astFile:   pkgInfo.Files[i],
			fullpath:  fullpath,
			undefined: undefined,
			derived:   derived,
			funcNames: f.funcNames,
		})
	}
	return files
}

type finder struct {
	program   *loader.Program
	pkgInfo   *loader.PackageInfo
	undefined []*ast.CallExpr
	derived   []*ast.CallExpr
	funcNames map[string]struct{}
}

func (f *finder) Visit(node ast.Node) (w ast.Visitor) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return f
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return f
	}
	def, ok := f.pkgInfo.Uses[fn]
	if !ok {
		f.undefined = append(f.undefined, call)
		return f
	}
	if _, ok := def.(*types.Builtin); ok {
		return f
	}
	file := f.program.Fset.File(def.Pos())
	if file == nil {
		// probably a cast, for example float64()
		return f
	}
	_, filename := filepath.Split(file.Name())
	if filename == derivedFilename {
		f.derived = append(f.derived, call)
		return f
	}
	f.funcNames[fn.Name] = struct{}{}
	return f
}

type call struct {
	Expr *ast.CallExpr
	Name string
	Args []types.Type
}

func newCall(pkgInfo *loader.PackageInfo, expr *ast.CallExpr) *call {
	fn, ok := expr.Fun.(*ast.Ident)
	if !ok {
		panic("unreachable, finder has already eliminated this option")
	}
	name := fn.Name
	typs := getInputTypes(pkgInfo, expr)
	return &call{expr, name, typs}
}

// argTypes returns the argument types of a function call.
func getInputTypes(pkgInfo *loader.PackageInfo, call *ast.CallExpr) []types.Type {
	typs := make([]types.Type, len(call.Args))
	for i, a := range call.Args {
		typs[i] = pkgInfo.TypeOf(a)
	}
	return typs
}

// HasUndefined returns whether the call has undefined arguments
func (c *call) HasUndefined() bool {
	for i := range c.Args {
		if c.Args[i] == nil {
			return true
		}
		if basic, ok := c.Args[i].(*types.Basic); ok {
			if basic.Kind() == types.Invalid {
				return true
			}
		}
	}
	return false
}
