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

	"golang.org/x/tools/go/loader"
)

// Program is ready to generate code for a whole program.
type Program interface {
	Generate() error
}

// Plugins is a collection of plugins,
// that given a list of paths becomes a Program.
type Plugins interface {
	Load(paths []string) (Program, error)
}

type plugins struct {
	plugins  []Plugin
	autoname bool
	dedup    bool
}

// NewPlugins returns a collection of plugins that is ready to generate code.
func NewPlugins(ps []Plugin, autoname bool, dedup bool) Plugins {
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

func (p *plugins) Load(paths []string) (Program, error) {
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

func union(this, that map[string]struct{}) map[string]struct{} {
	for k := range that {
		this[k] = struct{}{}
	}
	return this
}

func newPackage(program *loader.Program, pkgInfo *loader.PackageInfo, plugins []Plugin, autoname, dedup bool) (*pkg, error) {
	fileInfos := newFileInfos(program, pkgInfo)
	fullpath := ""
	if len(fileInfos) > 0 {
		abs, err := filepath.Abs(fileInfos[0].fullpath)
		if err != nil {
			return nil, err
		}
		fullpath = filepath.Dir(abs)
	}
	reserved := make(map[string]struct{})
	for _, fileFuncs := range fileInfos {
		reserved = union(reserved, fileFuncs.funcNames)
	}

	printer := newPrinter(pkgInfo.Pkg.Name())
	qual := newQualifier(printer, pkgInfo.Pkg)
	typesmaps := make(map[string]TypesMap, len(plugins))
	deps := make(map[string]Dependency, len(plugins))
	for _, plugin := range plugins {
		tm := newTypesMap(qual, plugin.GetPrefix(), reserved, autoname, dedup)
		deps[plugin.Name()] = tm
		typesmaps[plugin.Name()] = tm
	}
	generators := make(map[string]Generator, len(plugins))
	for _, plugin := range plugins {
		generators[plugin.Name()] = plugin.New(typesmaps[plugin.Name()], printer, deps)
	}
	pkg := &pkg{pkgInfo, plugins, generators, printer, nil, fullpath}
	for _, fileInfo := range fileInfos {

		changed := false
		calls := append(fileInfo.undefined, fileInfo.derived...)
		for _, call := range calls {
			// log.Printf("call: %v", call.Name)
			if call.HasUndefined() {
				// Only functions that are supported by a code generator plugin should be added to undefined.
				// Otherwise things like functions like casts will get into the undefined loop and result in
				// an error, where a function could not be generated.
				for _, plugin := range plugins {
					if strings.HasPrefix(call.Name, plugin.GetPrefix()) {
						pkg.undefined = append(pkg.undefined, call.Expr)
						break
					}
				}
				continue
			}
			name, err := pkg.Add(call)
			if err != nil {
				return nil, err
			}
			if len(name) == 0 {
				// this call did not match any prefixes of any code generator and is undefined.
				continue
			}
			if name != call.Name {
				if !autoname && !dedup {
					panic("unreachable: function names cannot be changed if it is not allowed by the user")
				}
				changed = true
				log.Printf("changing function call name from %s to %s", call.Name, name)
				call.Expr.Fun = ast.NewIdent(name)
			}
		}

		if changed {
			info, err := os.Stat(fileInfo.fullpath)
			if err != nil {
				return nil, fmt.Errorf("stat %s: %v", fileInfo.fullpath, err)
			}
			f, err := os.OpenFile(fileInfo.fullpath, os.O_WRONLY, info.Mode())
			if err != nil {
				return nil, fmt.Errorf("opening %s: %v", fileInfo.fullpath, err)
			}
			defer f.Close()
			if err := format.Node(f, program.Fset, fileInfo.astFile); err != nil {
				return nil, fmt.Errorf("formatting %s: %v", fileInfo.fullpath, err)
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
	fullpath   string
}

func (pkg *pkg) Add(call *call) (string, error) {
	for _, p := range pkg.plugins {
		if !strings.HasPrefix(call.Name, p.GetPrefix()) {
			continue
		}
		generator := pkg.generators[p.Name()]
		name, err := generator.Add(call.Name, call.Args)
		if err != nil {
			return "", fmt.Errorf("Add Error: %s: %v", p.Name(), err)
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

func (pkg *pkg) HasContent() bool {
	return pkg.printer.HasContent()
}

func (pkg *pkg) Filename() string {
	return filepath.Join(pkg.fullpath, derivedFilename)
}

func (pkg *pkg) Print() error {
	f, err := os.Create(pkg.Filename())
	if err != nil {
		return err
	}
	if _, err := pkg.printer.WriteTo(f); err != nil {
		return err
	}
	return f.Close()
}

func (pkg *pkg) Delete() error {
	filename := pkg.Filename()
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat %s: %v", filename, err)
	}
	return os.Remove(filename)
}

func (pkg *pkg) Generate() (bool, error) {
	generated := false
	for !pkg.Done() {
		for _, plugin := range pkg.plugins {
			g := pkg.generators[plugin.Name()]
			for _, typs := range g.ToGenerate() {
				if err := g.Generate(typs); err != nil {
					return false, fmt.Errorf("Generator Error: %s:%v", plugin.Name(), err.Error())
				}
				generated = true
			}
		}
	}
	return generated, nil
}

func (pg *program) Generate() error {
	pkgInfos := pg.program.InitialPackages()

	// sort.Slice(pkgInfos, func(i, j int) bool {
	// 	return pkgInfos[i].String() < pkgInfos[j].String()
	// })
	for i := range pkgInfos {
		if err := pg.generatePackage(pkgInfos[i]); err != nil {
			return err
		}
	}
	return nil
}

func (pg *program) generatePackage(pkgInfo *loader.PackageInfo) error {
	path := pkgInfo.Pkg.Path()
	// ss := make([]string, len(pkgInfo.Files))
	// for i := range pkgInfo.Files {
	// 	ss[i] = pg.program.Fset.File(pkgInfo.Files[i].Pos()).Name()
	// }
	// log.Printf("package: %s, files %d: %s", path, len(pkgInfo.Files), strings.Join(ss, ", "))
	generated := true
	var undefined string
	thisprogram := pg.program
	for generated {
		pkgGen, err := newPackage(thisprogram, pkgInfo, pg.plugins, pg.autoname, pg.dedup)
		if err != nil {
			return err
		}

		us := make([]string, len(pkgGen.undefined))
		for i, u := range pkgGen.undefined {
			us[i] = types.ExprString(u)
		}
		sort.Strings(us)

		for _, u := range us {
			log.Printf("could not yet generate: %s", u)
		}

		generated, err = pkgGen.Generate()
		if err != nil {
			return err
		}

		if pkgGen.HasContent() {
			if err := pkgGen.Print(); err != nil {
				return err
			}
		} else {
			// When the file has no content it should be removed.
			if err := pkgGen.Delete(); err != nil {
				return err
			}
		}

		if len(us) == 0 {
			return nil
		}

		newundefined := strings.Join(us, ";")
		if newundefined == undefined {
			break
		}
		undefined = newundefined

		// reload path with newly generated code, with the hope that some types are now inferable.
		thisprogram, err = load(path)
		if err != nil {
			return err
		}
		pkgInfo = thisprogram.Package(path)
	}

	if len(undefined) > 0 && !generated {
		return fmt.Errorf("cannot generate: %s", undefined)
	}
	return nil
}
