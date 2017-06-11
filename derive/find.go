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

func findUndefinedOrDerivedFuncs(program *loader.Program, pkgInfo *loader.PackageInfo, file *ast.File) []*Call {
	f := &finder{program, pkgInfo, nil, nil}
	for _, d := range file.Decls {
		ast.Walk(f, d)
	}
	callExprs := append(f.undefined, f.derived...)
	calls := make([]*Call, len(callExprs))
	for i := range callExprs {
		calls[i] = newCall(pkgInfo, callExprs[i])
	}
	return calls
}

type finder struct {
	program   *loader.Program
	pkgInfo   *loader.PackageInfo
	undefined []*ast.CallExpr
	derived   []*ast.CallExpr
}

func (this *finder) Visit(node ast.Node) (w ast.Visitor) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return this
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return this
	}
	def, ok := this.pkgInfo.Uses[fn]
	if !ok {
		this.undefined = append(this.undefined, call)
		return this
	}
	file := this.program.Fset.File(def.Pos())
	if file == nil {
		// probably a builtin function, for example panic.
		return this
	}
	_, filename := filepath.Split(file.Name())
	if filename == derivedFilename {
		this.derived = append(this.derived, call)
	}
	return this
}

type Call struct {
	Expr *ast.CallExpr
	Name string
	Args []types.Type
}

func newCall(pkgInfo *loader.PackageInfo, expr *ast.CallExpr) *Call {
	fn, ok := expr.Fun.(*ast.Ident)
	if !ok {
		panic("unreachable, finder has already eliminated this option")
	}
	name := fn.Name
	typs := getInputTypes(pkgInfo, expr)
	return &Call{expr, name, typs}
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
func (this *Call) HasUndefined() bool {
	for i := range this.Args {
		if this.Args[i] == nil {
			return true
		}
	}
	return false
}
