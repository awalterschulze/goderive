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
	"fmt"
	"go/ast"
	"go/format"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

type programGenerator struct {
	plugins  []Plugin
	autoname bool
	dedup    bool
	program  *loader.Program
}

func NewGenerator(plugins []Plugin, paths []string, autoname bool, dedup bool) (*programGenerator, error) {
	program, err := load(paths...)
	if err != nil {
		return nil, err
	}
	return &programGenerator{
		plugins:  plugins,
		autoname: autoname,
		dedup:    dedup,
		program:  program,
	}, nil
}

func (p *programGenerator) NewPackage(pkgInfo *loader.PackageInfo) *packageGenerator {
	printer := newPrinter(pkgInfo.Pkg.Name())
	qual := newQualifier(printer, pkgInfo.Pkg)
	typesmaps := make(map[string]TypesMap, len(p.plugins))
	deps := make(map[string]Dependency, len(p.plugins))
	for _, plugin := range p.plugins {
		tm := newTypesMap(qual, plugin.GetPrefix(), p.autoname, p.dedup)
		deps[plugin.Name()] = tm
		typesmaps[plugin.Name()] = tm
	}

	generators := make(map[string]Generator, len(p.plugins))
	for _, plugin := range p.plugins {
		generators[plugin.Name()] = plugin.New(typesmaps[plugin.Name()], printer, deps)
	}
	return &packageGenerator{p.plugins, generators, printer}
}

type packageGenerator struct {
	plugins    []Plugin
	generators map[string]Generator
	printer    Printer
}

func (g *packageGenerator) Add(call *Call) (string, error) {
	for _, p := range g.plugins {
		if !strings.HasPrefix(call.Name, p.GetPrefix()) {
			continue
		}
		generator := g.generators[p.Name()]
		name, err := generator.Add(call.Name, call.Args)
		if err != nil {
			return "", fmt.Errorf("%s: %v", p.Name(), err)
		}
		return name, nil
	}
	return "", nil
}

func (pg *programGenerator) Generate() error {
	for _, pkgInfo := range pg.program.InitialPackages() {

		path := pkgInfo.Pkg.Path()
		ungenerated := -1
		for ungenerated != 0 {
			thisprogram := pg.program
			if ungenerated > 0 {
				// reload path with newly generated code, with the hope that some types are now inferable.
				thisprogram, err := load(path)
				if err != nil {
					return err
				}
				pkgInfo = thisprogram.Package(path)
			}
			pkgGen := pg.NewPackage(pkgInfo)

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

				calls := findUndefinedOrDerivedFuncs(thisprogram, pkgInfo, file)
				changed := false
				for _, call := range calls {
					if call.HasUndefined() {
						notgenerated = append(notgenerated, call.Expr)
						continue
					}
					name, err := pkgGen.Add(call)
					if err != nil {
						return err
					}
					generated := len(name) > 0
					if generated {
						if name != call.Name {
							if !pg.autoname && !pg.dedup {
								panic("unreachable: function names cannot be changed if it is not allowed by the user")
							}
							changed = true
							log.Printf("changing function call name from %s to %s", call.Name, name)
							call.Expr.Fun = ast.NewIdent(name)
						}
					} else {
						notgenerated = append(notgenerated, call.Expr)
					}
				}
				if changed {
					info, err := os.Stat(fullpath)
					if err != nil {
						return fmt.Errorf("stat %s: %v", fullpath, err)
					}
					f, err := os.OpenFile(fullpath, os.O_WRONLY, info.Mode())
					if err != nil {
						return fmt.Errorf("opening %s: %v", fullpath, err)
					}
					if err := format.Node(f, thisprogram.Fset, file); err != nil {
						return fmt.Errorf("formatting %s: %v", fullpath, err)
					}
				}
			}

			if len(notgenerated) > 0 {
				if ungenerated == len(notgenerated) {
					for _, c := range notgenerated {
						return fmt.Errorf("cannot generate %s", types.ExprString(c))
					}
				}
			}
			ungenerated = len(notgenerated)
			for _, c := range notgenerated {
				log.Printf("could not yet generate: %s", types.ExprString(c))
			}

			alldone := false
			for !alldone {
				for _, plugin := range pkgGen.plugins {
					if err := pkgGen.generators[plugin.Name()].Generate(); err != nil {
						return fmt.Errorf(plugin.Name() + ":" + err.Error())
					}
				}
				alldone = func() bool {
					for _, g := range pkgGen.generators {
						if !g.Done() {
							return false
						}
					}
					return true
				}()
			}

			if pkgGen.printer.HasContent() {
				pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkgInfo.Pkg.Path())
				f, err := os.Create(filepath.Join(pkgpath, derivedFilename))
				if err != nil {
					return err
				}
				if err := pkgGen.printer.WriteTo(f); err != nil {
					return err
				}
				f.Close()
			}
		}
	}
	return nil
}
