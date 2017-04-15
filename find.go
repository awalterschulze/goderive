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

package main

import (
	"go/ast"
	"path/filepath"

	"golang.org/x/tools/go/loader"
)

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

func findUndefinedOrDerivedFuncs(program *loader.Program, pkgInfo *loader.PackageInfo, file *ast.File) []*ast.CallExpr {
	f := &finder{program, pkgInfo, nil, nil}
	for _, d := range file.Decls {
		ast.Walk(f, d)
	}
	return append(f.undefined, f.derived...)
}
