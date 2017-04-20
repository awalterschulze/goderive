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

	"golang.org/x/tools/go/loader"

	"github.com/kisielk/gotool"
)

const derivedFilename = "derived.gen.go"

type Generator interface {
	TypesMap
	Add(pkgInfo *loader.PackageInfo, call *ast.CallExpr) (bool, error)
	Generate() error
}

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
			keysTypesMap := newTypesMap(qual, *keysPrefix)
			sortedTypesMap := newTypesMap(qual, *sortedPrefix)
			compareTypesMap := newTypesMap(qual, *comparePrefix)
			fmapTypesMap := newTypesMap(qual, *fmapPrefix)
			joinTypesMap := newTypesMap(qual, *joinPrefix)

			generators := []Generator{
				newEqual(equalTypesMap, qual, *equalPrefix, p),
				newKeys(keysTypesMap, qual, *keysPrefix, p),
				newCompare(compareTypesMap, qual, *comparePrefix, p, keysTypesMap, sortedTypesMap),
				newSorted(sortedTypesMap, qual, *sortedPrefix, p, compareTypesMap),
				newFmap(fmapTypesMap, qual, *fmapPrefix, p),
				newJoin(joinTypesMap, qual, *joinPrefix, p),
			}

			var err error
			for _, call := range calls {
				generated := false
				for _, gen := range generators {
					generated, err = gen.Add(pkgInfo, call)
					if err != nil {
						log.Fatal("equal:" + err.Error())
					} else if generated {
						break
					}
				}
				if !generated {
					notgenerated = append(notgenerated, call)
				}
			}

			if len(notgenerated) > 0 {
				if ungenerated == len(notgenerated) {
					for _, c := range notgenerated {
						log.Fatalf("cannot generate %s", types.ExprString(c))
					}
				}
			}
			ungenerated = len(notgenerated)
			for _, c := range notgenerated {
				log.Printf("could not yet generate: %s", types.ExprString(c))
			}

			alldone := false
			for !alldone {
				for _, gen := range generators {
					if err := gen.Generate(); err != nil {
						log.Fatal(err)
					}
				}
				alldone = true
				for _, gen := range generators {
					if !gen.Done() {
						alldone = false
						break
					}
				}
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
}
