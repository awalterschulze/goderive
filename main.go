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
	"flag"
	"go/ast"
	"log"
	"os"
	"path/filepath"

	"github.com/kisielk/gotool"
)

const derivedFilename = "derived.gen.go"

func main() {
	log.SetFlags(0)
	flag.Parse()
	paths := gotool.ImportPaths(flag.Args())
	program, err := load(paths...)
	if err != nil {
		log.Fatal(err) // load error
	}
	for _, pkgInfo := range program.InitialPackages() {

		var calls []*ast.CallExpr
		for _, file := range pkgInfo.Files {
			astFile := program.Fset.File(file.Pos())
			if astFile == nil {
				continue
			}
			_, fname := filepath.Split(astFile.Name())
			if fname == derivedFilename {
				continue
			}
			newcalls := findUndefinedOrDerivedFuncs(program, pkgInfo, file)
			calls = append(newcalls, calls...)
		}

		p := newPrinter(pkgInfo.Pkg.Name())

		if err := generateEqual(p, pkgInfo, *equalPrefix, false, calls); err != nil {
			log.Fatal(err)
		}
		if err := generateSortedKeys(p, pkgInfo, *sortedKeysPrefix, false, calls); err != nil {
			log.Fatal(err)
		}
		if err := generateCompare(p, pkgInfo, *comparePrefix, false, calls); err != nil {
			log.Fatal(err)
		}

		if p.HasContent() {
			pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkgInfo.Pkg.Path())
			f, err := os.Create(filepath.Join(pkgpath, derivedFilename))
			if err != nil {
				log.Fatal(err)
			}
			if err := p.WriteTo(f); err != nil {
				panic(err)
			}
			f.Close()
		}

	}
}
