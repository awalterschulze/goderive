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
	"go/format"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/awalterschulze/goderive/derive"
	"github.com/awalterschulze/goderive/plugin/compare"
	"github.com/awalterschulze/goderive/plugin/equal"
	"github.com/awalterschulze/goderive/plugin/fmap"
	"github.com/awalterschulze/goderive/plugin/join"
	"github.com/awalterschulze/goderive/plugin/keys"
	"github.com/awalterschulze/goderive/plugin/sorted"
	"github.com/kisielk/gotool"
)

var autoname = flag.Bool("autoname", false, "rename functions that are conflicting with other functions")
var dedup = flag.Bool("dedup", false, "rename functions to functions that are duplicates")

const derivedFilename = "derived.gen.go"

type Generator interface {
	derive.TypesMap
	Add(name string, typs []types.Type) (string, error)
	Name() string
	Generate() error
}

func main() {
	log.SetFlags(0)
	flag.Parse()
	paths := gotool.ImportPaths(flag.Args())

	program, err := derive.Load(paths...)
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
				thisprogram, err = derive.Load(path)
				if err != nil {
					log.Fatal(err)
				}
				pkgInfo = thisprogram.Package(path)
			}

			p := derive.NewPrinter(pkgInfo.Pkg.Name())
			quals := derive.NewQualifier(p, pkgInfo.Pkg)

			equalTypesMap := derive.NewTypesMap(quals, *equal.Prefix, *autoname, *dedup)
			keysTypesMap := derive.NewTypesMap(quals, *keys.Prefix, *autoname, *dedup)
			sortedTypesMap := derive.NewTypesMap(quals, *sorted.Prefix, *autoname, *dedup)
			compareTypesMap := derive.NewTypesMap(quals, *compare.Prefix, *autoname, *dedup)
			fmapTypesMap := derive.NewTypesMap(quals, *fmap.Prefix, *autoname, *dedup)
			joinTypesMap := derive.NewTypesMap(quals, *join.Prefix, *autoname, *dedup)

			generators := []Generator{
				equal.New(equalTypesMap, p),
				keys.New(keysTypesMap, p),
				compare.New(compareTypesMap, p, keysTypesMap, sortedTypesMap),
				sorted.New(sortedTypesMap, p, compareTypesMap),
				fmap.New(fmapTypesMap, p),
				join.New(joinTypesMap, p),
			}

			var notgenerated []*ast.CallExpr
			for _, file := range pkgInfo.Files {
				astFile := thisprogram.Fset.File(file.Pos())
				if astFile == nil {
					continue
				}
				fullpath := astFile.Name()
				_, fname := filepath.Split(fullpath)
				if fname == derivedFilename {
					continue
				}

				calls := derive.FindUndefinedOrDerivedFuncs(thisprogram, pkgInfo, file)
				changed := false
				for _, call := range calls {
					if call.HasUndefined() {
						notgenerated = append(notgenerated, call.Expr)
						continue
					}
					generated := func() bool {
						for _, gen := range generators {
							if !strings.HasPrefix(call.Name, gen.Prefix()) {
								continue
							}
							name, err := gen.Add(call.Name, call.Args)
							if err != nil {
								log.Fatalf("%s: %v", gen.Name(), err)
							}
							if name != call.Name {
								if !*autoname && !*dedup {
									panic("unreachable: function names cannot be changed if it is not allowed by the user")
								}
								changed = true
								log.Printf("changing function call name from %s to %s", call.Name, name)
								call.Expr.Fun = ast.NewIdent(name)
							}
							return true
						}
						return false
					}()
					if !generated {
						notgenerated = append(notgenerated, call.Expr)
					}
				}
				if changed {
					info, err := os.Stat(fullpath)
					if err != nil {
						log.Fatalf("stat %s: %v", fullpath, err)
					}
					f, err := os.OpenFile(fullpath, os.O_WRONLY, info.Mode())
					if err != nil {
						log.Fatalf("opening %s: %v", fullpath, err)
					}
					if err := format.Node(f, thisprogram.Fset, file); err != nil {
						log.Fatalf("formatting %s: %v", fullpath, err)
					}
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
						log.Fatal(gen.Name() + ":" + err.Error())
					}
				}
				alldone = func() bool {
					for _, gen := range generators {
						if !gen.Done() {
							return false
						}
					}
					return true
				}()
			}

			if p.HasContent() {
				pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkgInfo.Pkg.Path())
				f, err := os.Create(filepath.Join(pkgpath, derivedFilename))
				if err != nil {
					log.Fatal(err)
				}
				if err := p.WriteTo(f); err != nil {
					log.Fatal(err)
				}
				f.Close()
			}
		}
	}
}
