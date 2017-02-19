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
	"fmt"
	"go/ast"
	"go/types"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

type finder struct {
	program *loader.Program
	pkgInfo *loader.PackageInfo
	calls   []*ast.CallExpr
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
		// undefined function.
		this.calls = append(this.calls, call)
		return this
	}
	file := this.program.Fset.File(def.Pos())
	if file == nil {
		// probably a builtin function, for example panic.
		return this
	}
	_, filename := filepath.Split(file.Name())
	if filename == derivedFilename {
		// derived function.
		this.calls = append(this.calls, call)
	}
	return this
}

func findUndefinedOrDerivedFuncs(program *loader.Program, pkgInfo *loader.PackageInfo, file *ast.File) []*ast.CallExpr {
	f := &finder{program, pkgInfo, nil}
	for _, d := range file.Decls {
		ast.Walk(f, d)
	}
	return f.calls
}

func findEqualFuncs(pkgInfo *loader.PackageInfo, calls []*ast.CallExpr) []types.Type {
	var typs []types.Type
	for _, call := range calls {
		fn, ok := call.Fun.(*ast.Ident)
		if !ok {
			continue
		}
		if !strings.HasPrefix(fn.Name, eqFuncPrefix) {
			continue
		}
		if len(call.Args) != 2 {
			fmt.Fprintf(os.Stderr, "%s does not have two arguments\n", fn.Name)
			continue
		}
		t0 := pkgInfo.TypeOf(call.Args[0])
		t1 := pkgInfo.TypeOf(call.Args[1])
		if !types.Identical(t0, t1) {
			fmt.Fprintf(os.Stderr, "%s has two arguments, but they are of different types %s != %s\n",
				fn.Name, t0, t1)
			continue
		}
		name := strings.TrimPrefix(fn.Name, eqFuncPrefix)
		qual := types.RelativeTo(pkgInfo.Pkg)
		typeStr := typeName(t0, qual)
		if typeStr != name {
			//TODO think about whether this is really necessary
			fmt.Fprintf(os.Stderr, "%s's suffix %s does not match the type %s\n",
				fn.Name, name, typeStr)
			continue
		}
		typs = append(typs, t0)
	}
	return typs
}
