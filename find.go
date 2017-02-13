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
	"strings"

	"golang.org/x/tools/go/loader"
)

func findTypesForFuncPrefix(pkgInfo *loader.PackageInfo, f *ast.File, funcPrefix string) []types.Type {
	var typs []types.Type
	for _, d := range f.Decls {
		finder := &findTypes{pkgInfo, funcPrefix, nil}
		ast.Walk(finder, d)
		typs = append(typs, finder.typs...)
	}
	return typs
}

type findTypes struct {
	pkgInfo    *loader.PackageInfo
	funcPrefix string
	typs       []types.Type
}

func (this *findTypes) Visit(node ast.Node) (w ast.Visitor) {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return this
	}
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return this
	}
	if !strings.HasPrefix(fn.Name, this.funcPrefix) {
		return this
	}
	if len(call.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s does not have two arguments\n", fn.Name)
		return this
	}
	t0 := this.pkgInfo.TypeOf(call.Args[0])
	t1 := this.pkgInfo.TypeOf(call.Args[1])
	if !types.Identical(t0, t1) {
		fmt.Fprintf(os.Stderr, "%s has two arguments, but they are of different types %s != %s\n",
			fn.Name, t0, t1)
		return this
	}
	name := strings.TrimPrefix(fn.Name, this.funcPrefix)
	qual := types.RelativeTo(this.pkgInfo.Pkg)
	typeStr := typeName(t0, qual)
	if typeStr != name {
		//TODO think about whether this is really necessary
		fmt.Fprintf(os.Stderr, "%s's suffix %s does not match the type %s\n",
			fn.Name, name, typeStr)
		return this
	}
	this.typs = append(this.typs, t0)
	return this
}
