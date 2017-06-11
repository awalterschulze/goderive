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
	"sort"
	"strings"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

type plugins struct {
	plugins  []Plugin
	autoname bool
	dedup    bool
}

func NewPlugins(ps []Plugin, autoname bool, dedup bool) *plugins {
	sortPlugins(ps)
	return &plugins{
		plugins:  ps,
		autoname: autoname,
		dedup:    dedup,
	}
}

// sortPlugins sorts plugins from biggest to smallest prefix to make sure than conflicts in prefixes are resolved.
// For example: derivSorted should generated a sorted function and not a sort function.
func sortPlugins(ps []Plugin) {
	sort.Slice(ps, func(i, j int) bool {
		if len(ps[i].GetPrefix()) == len(ps[j].GetPrefix()) {
			return ps[i].GetPrefix() > ps[j].GetPrefix()
		}
		return len(ps[i].GetPrefix()) > len(ps[j].GetPrefix())
	})
}

type program struct {
	plugins  []Plugin
	autoname bool
	dedup    bool
	program  *loader.Program
}

func (p *plugins) Load(paths []string) (*program, error) {
	loaded, err := load(paths...)
	if err != nil {
		return nil, err
	}
	return &program{
		plugins:  p.plugins,
		autoname: p.autoname,
		dedup:    p.dedup,
		program:  loaded,
	}, nil
}

func (p *program) NewPackage(pkgInfo *loader.PackageInfo) (*pkg, error) {
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
	pkg := &pkg{pkgInfo, p.plugins, generators, printer, nil}
	for _, file := range pkgInfo.Files {
		astFile := p.program.Fset.File(file.Pos())
		if astFile == nil {
			continue
		}
		fullpath := astFile.Name()
		_, fname := filepath.Split(fullpath)
		if fname == derivedFilename {
			continue
		}

		calls := findUndefinedOrDerivedFuncs(p.program, pkgInfo, file)
		changed := false
		for _, call := range calls {
			if call.HasUndefined() {
				pkg.undefined = append(pkg.undefined, call.Expr)
				continue
			}
			name, err := pkg.Add(call)
			if err != nil {
				return nil, err
			}
			generated := len(name) > 0
			if generated {
				if name != call.Name {
					if !p.autoname && !p.dedup {
						panic("unreachable: function names cannot be changed if it is not allowed by the user")
					}
					changed = true
					log.Printf("changing function call name from %s to %s", call.Name, name)
					call.Expr.Fun = ast.NewIdent(name)
				}
			} else {
				pkg.undefined = append(pkg.undefined, call.Expr)
			}
		}
		if changed {
			info, err := os.Stat(fullpath)
			if err != nil {
				return nil, fmt.Errorf("stat %s: %v", fullpath, err)
			}
			f, err := os.OpenFile(fullpath, os.O_WRONLY, info.Mode())
			if err != nil {
				return nil, fmt.Errorf("opening %s: %v", fullpath, err)
			}
			if err := format.Node(f, p.program.Fset, file); err != nil {
				return nil, fmt.Errorf("formatting %s: %v", fullpath, err)
			}
		}
	}
	return pkg, nil
}

type pkg struct {
	info       *loader.PackageInfo
	plugins    []Plugin
	generators map[string]Generator
	printer    Printer
	undefined  []*ast.CallExpr
}

func (pkg *pkg) Add(call *Call) (string, error) {
	for _, p := range pkg.plugins {
		if !strings.HasPrefix(call.Name, p.GetPrefix()) {
			continue
		}
		generator := pkg.generators[p.Name()]
		name, err := generator.Add(call.Name, call.Args)
		if err != nil {
			return "", fmt.Errorf("%s: %v", p.Name(), err)
		}
		return name, nil
	}
	return "", nil
}

func (pkg *pkg) Done() bool {
	for _, g := range pkg.generators {
		if !g.Done() {
			return false
		}
	}
	return true
}

func (pkg *pkg) Print() error {
	if !pkg.printer.HasContent() {
		return nil
	}
	pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkg.info.Pkg.Path())
	f, err := os.Create(filepath.Join(pkgpath, derivedFilename))
	if err != nil {
		return err
	}
	if err := pkg.printer.WriteTo(f); err != nil {
		return err
	}
	return f.Close()
}

func (pkg *pkg) Generate() error {
	for !pkg.Done() {
		for _, plugin := range pkg.plugins {
			if err := pkg.generators[plugin.Name()].Generate(); err != nil {
				return fmt.Errorf(plugin.Name() + ":" + err.Error())
			}
		}
	}
	return nil
}

func (pg *program) Generate() error {
	for _, pkgInfo := range pg.program.InitialPackages() {
		if err := pg.generatePackage(pkgInfo); err != nil {
			return err
		}
	}
	return nil
}

func (pg *program) generatePackage(pkgInfo *loader.PackageInfo) error {
	path := pkgInfo.Pkg.Path()
	prevNumUndefined := -1
	for {
		pkgGen, err := pg.NewPackage(pkgInfo)
		if err != nil {
			return err
		}
		undefined := pkgGen.undefined

		if len(undefined) > 0 {
			if prevNumUndefined == len(undefined) {
				ss := make([]string, len(undefined))
				for i, c := range undefined {
					ss[i] = types.ExprString(c)
				}
				return fmt.Errorf("cannot generate: %s", strings.Join(ss, ","))
			}
		}
		prevNumUndefined = len(undefined)
		for _, c := range undefined {
			log.Printf("could not yet generate: %s", types.ExprString(c))
		}

		if err := pkgGen.Generate(); err != nil {
			return err
		}

		if err := pkgGen.Print(); err != nil {
			return err
		}

		if len(undefined) == 0 {
			return nil
		}

		// reload path with newly generated code, with the hope that some types are now inferable.
		thisprogram, err := load(path)
		if err != nil {
			return err
		}
		pkgInfo = thisprogram.Package(path)
	}
}
