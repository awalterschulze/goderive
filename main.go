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
	"go/types"
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
		log.Fatal(err)
	}
	for _, pkgInfo := range program.InitialPackages() {
		path := pkgInfo.Pkg.Path()
		ungenerated := -1
		for ungenerated != 0 {
			thisprogram := program
			if ungenerated > 0 {
				// reload path with newly generated code, with the hope that some types are now inferable.
				thisprogram, err = load(path)
				if err != nil {
					log.Fatal(err)
				}
				pkgInfo = thisprogram.Package(path)
			}
			var notgenerated []*ast.CallExpr
			var calls []*ast.CallExpr
			for _, file := range pkgInfo.Files {
				astFile := thisprogram.Fset.File(file.Pos())
				if astFile == nil {
					continue
				}
				_, fname := filepath.Split(astFile.Name())
				if fname == derivedFilename {
					continue
				}
				newcalls := findUndefinedOrDerivedFuncs(thisprogram, pkgInfo, file)
				calls = append(newcalls, calls...)
			}
			qual := types.RelativeTo(pkgInfo.Pkg)

			p := newPrinter(pkgInfo.Pkg.Name())

			equalTypesMap := newTypesMap(qual, *equalPrefix)
			sortedKeysTypesMap := newTypesMap(qual, *sortedKeysPrefix)
			compareTypesMap := newTypesMap(qual, *comparePrefix)

			equal, err := newEqual(p, qual, equalTypesMap)
			if err != nil {
				log.Fatal(err)
			}
			sortedKeys, err := newSortedKeys(p, qual, sortedKeysTypesMap, compareTypesMap)
			if err != nil {
				log.Fatal(err)
			}
			compare, err := newCompare(p, qual, compareTypesMap, sortedKeysTypesMap)
			if err != nil {
				log.Fatal(err)
			}

			alldone := false
			for !alldone {
				for _, call := range calls {
					if generated, err := equal.Generate(pkgInfo, *equalPrefix, call); err != nil {
						log.Fatal(err)
					} else if generated {
						continue
					}
					if generated, err := sortedKeys.Generate(pkgInfo, *sortedKeysPrefix, call); err != nil {
						log.Fatal(err)
					} else if generated {
						continue
					}
					if generated, err := compare.Generate(pkgInfo, *comparePrefix, call); err != nil {
						log.Fatal(err)
					} else if generated {
						continue
					}
					notgenerated = append(notgenerated, call)
				}
				alldone = equal.Done() && sortedKeys.Done() && compare.Done()
			}

			if len(notgenerated) > 0 {
				if ungenerated == len(notgenerated) {
					log.Fatalf("cannot generate %v", notgenerated)
				}
			}
			ungenerated = len(notgenerated)

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
}
