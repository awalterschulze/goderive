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
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kisielk/gotool"
)

func main() {
	flag.Parse()
	paths := gotool.ImportPaths(flag.Args())
	program, err := load(paths...)
	if err != nil {
		log.Fatal(err) // load error
	}
	derivedFilename := "derived.gen.go"
	for _, pkgInfo := range program.InitialPackages() {
		pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkgInfo.Pkg.Path())

		qual := types.RelativeTo(pkgInfo.Pkg)
		var typs []types.Type
		for i, f := range pkgInfo.Files {
			gotFile := program.Fset.File(f.Pos())
			if gotFile == nil {
				continue
			}
			_, fname := filepath.Split(gotFile.Name())
			if fname == derivedFilename {
				continue
			}
			newtyps := findTypesForFuncPrefix(pkgInfo, pkgInfo.Files[i], eqFuncPrefix)
			typs = append(typs, newtyps...)
		}
		if len(typs) == 0 {
			continue
		}

		f, err := os.Create(filepath.Join(pkgpath, derivedFilename))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		generate(f, qual, pkgInfo.Pkg.Name(), typs)
	}
}

func generate(w io.Writer, qual types.Qualifier, pkgName string, typs []types.Type) {
	p := newPrinter()
	m := newTypesMap(qual)
	eq := newEqual(p, m, qual)
	for _, typ := range typs {
		m.Set(typ, false)
	}
	for _, typ := range typs {
		m.Set(typ, false)
		eq.gen(typ)
		m.Set(typ, true)
	}
	for _, typ := range m.List() {
		if m.Get(typ) {
			continue
		}
		eq.gen(typ)
		m.Set(typ, true)
	}
	p.Flush(pkgName, w)
}
