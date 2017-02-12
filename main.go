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
	"fmt"
	"go/ast"
	"go/parser"
	"go/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kisielk/gotool"
	"golang.org/x/tools/go/loader"
)

func isComparable(tt types.Type) bool {
	t := tt.Underlying()
	switch typ := t.(type) {
	case *types.Basic:
		return typ.Kind() != types.UntypedNil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			ft := f.Type()
			if !isComparable(ft) {
				return false
			}
		}
		return true
	case *types.Array:
		return isComparable(typ.Elem())
	}
	return false
}

type p struct {
	w              io.Writer
	qual           types.Qualifier
	generatedFuncs map[string]bool
	funcsToType    map[string]types.Type
}

func (p *p) And(format string, a ...interface{}) {
	fmt.Fprintf(p.w, "&&\n")
	fmt.Fprintf(p.w, format, a...)
}

func (p *p) P(format string, a ...interface{}) {
	fmt.Fprintf(p.w, format, a...)
}

func load(paths ...string) (*loader.Program, error) {
	conf := loader.Config{
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	rest, err := conf.FromArgs(paths, true)
	if err != nil {
		return nil, fmt.Errorf("could not parse arguments: %s", err)
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("unhandled extra arguments: %v", rest)
	}
	return conf.Load()
}

func main() {
	flag.Parse()
	paths := gotool.ImportPaths(flag.Args())
	program, err := load(paths...)
	eqFilename := "equal.gen.go"
	if err != nil {
		log.Fatal(err) // load error
	}
	for _, pkgInfo := range program.InitialPackages() {
		pkgpath := filepath.Join(filepath.Join(gotool.DefaultContext.BuildContext.GOPATH, "src"), pkgInfo.Pkg.Path())

		qual := types.RelativeTo(pkgInfo.Pkg)
		var typs []types.Type
		for i, f := range pkgInfo.Files {
			_, fname := filepath.Split(program.Fset.File(f.Pos()).Name())
			if fname == eqFilename {
				continue
			}
			newtyps := findObjectsUsingFunction(pkgInfo, pkgInfo.Files[i], eqFuncPrefix)
			typs = append(typs, newtyps...)
		}
		if len(typs) == 0 {
			continue
		}

		f, err := os.Create(filepath.Join(pkgpath, eqFilename))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		generate(f, qual, pkgInfo.Pkg.Name(), typs)
	}
}

func generate(w io.Writer, qual types.Qualifier, pkgName string, typs []types.Type) {
	p := &p{
		w:              w,
		qual:           qual,
		generatedFuncs: make(map[string]bool),
		funcsToType:    make(map[string]types.Type),
	}
	p.P("package %s\n", pkgName)
	p.P("\n")
	p.P("import (\n")
	p.P("\t\"bytes\"\n")
	p.P(")\n")
	p.P("\n")
	p.P("var _ = bytes.MinRead\n")
	p.P("\n")
	for i := range typs {
		p.newFunc(typs[i])
		p.genFunc(typs[i])
		p.funcGenerated(typs[i])
	}
	for name, generated := range p.generatedFuncs {
		if generated {
			continue
		}
		p.genFunc(p.funcsToType[name])
	}
}

const eqFuncPrefix = "derivEqual"

func funcName(typ types.Type, qual types.Qualifier) string {
	return eqFuncPrefix + typeName(typ, qual)
}

func (p *p) newFunc(typ types.Type) {
	fName := funcName(typ, p.qual)
	p.funcsToType[fName] = typ
	p.generatedFuncs[fName] = false
}

func typeName(typ types.Type, qual types.Qualifier) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return "PtrTo" + typeName(t.Elem(), qual)
	case *types.Array:
		return "ArrayOf" + typeName(t.Elem(), qual)
	case *types.Slice:
		return "SliceOf" + typeName(t.Elem(), qual)
	case *types.Map:
		return "MapOf" + typeName(t.Key(), qual) + "To" + typeName(t.Elem(), qual)
	}
	return types.TypeString(typ, qual)
}

func (p *p) funcGenerated(typ types.Type) {
	fName := funcName(typ, p.qual)
	p.generatedFuncs[fName] = true
}

func (p *p) genFunc(typ types.Type) {
	p.funcGenerated(typ)
	typeStr := types.TypeString(typ, p.qual)
	p.P("func %s(this, that %s) bool {\n", funcName(typ, p.qual), typeStr)
	switch ttyp := typ.(type) {
	case *types.Pointer:
		ref := ttyp.Elem()
		switch tttyp := ref.Underlying().(type) {
		case *types.Struct:
			p.P("return (%s == nil && %s == nil) ||\n", "this", "that")
			p.P("(%s != nil && %s != nil)", "this", "that")
			numFields := tttyp.NumFields()
			for i := 0; i < numFields; i++ {
				field := tttyp.Field(i)
				fieldType := field.Type()
				fieldName := field.Name()
				this := "this." + fieldName
				that := "that." + fieldName
				fieldStr, err := p.genField(this, that, fieldType)
				if err != nil {
					fmt.Fprintf(os.Stderr, err.Error())
					return
				}
				p.And(fieldStr)
			}
		default:
			fmt.Fprintf(os.Stderr, "unsupported: pointer is not a named struct, but %#v\n", ref)
			return
		}
	case *types.Slice:
		p.P("if this == nil {\n")
		p.P("if that == nil {\n")
		p.P("return true\n")
		p.P("} else {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("} else if that == nil {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("if len(this) != len(that) {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("for i := 0; i < len(this); i++ {\n")
		eqStr, err := p.genField("this[i]", "that[i]", ttyp.Elem())
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		p.P("if !(%s) {\n", eqStr)
		p.P("return false\n")
		p.P("}\n")
		p.P("}\n")
		p.P("return true\n")
	// TODO case *types.Array:

	case *types.Map:
		p.P("if this == nil {\n")
		p.P("if that == nil {\n")
		p.P("return true\n")
		p.P("} else {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("} else if that == nil {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("if len(this) != len(that) {\n")
		p.P("return false\n")
		p.P("}\n")
		p.P("for k, v := range this {\n")
		p.P("thatv, ok := that[k]\n")
		p.P("if !ok {\n")
		p.P("return false\n")
		p.P("}\n")
		eqStr, err := p.genField("v", "thatv", ttyp.Elem())
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		p.P("if !(%s) {\n", eqStr)
		p.P("return false\n")
		p.P("}\n")
		p.P("}\n")
		p.P("return true\n")
	default:
		fmt.Fprintf(os.Stderr, "unsupported type: %#v", typ)
		return
	}
	p.P("\n}\n")
}

func (p *p) genField(this, that string, fieldType types.Type) (string, error) {
	if isComparable(fieldType) {
		return fmt.Sprintf("%s == %s", this, that), nil
	}
	switch typ := fieldType.(type) {
	case *types.Pointer:
		ref := typ.Elem()
		if _, ok := ref.(*types.Named); ok {
			return fmt.Sprintf("%s.Equal(%s)", this, that), nil
		}
		eqStr, err := p.genField("*"+this, "*"+that, ref)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("((%[1]s == nil && %[2]s == nil) || (%[1]s != nil && %[2]s != nil && %[3]s))", this, that, eqStr), nil
	// TODO case *types.Array:
	// 	p.newFunc(typ)
	// 	return fmt.Sprintf("%s(%s, %s)", funcName(typ, p.qual), this, that), nil
	case *types.Slice:
		if b, ok := typ.Elem().(*types.Basic); ok && b.Kind() == types.Byte {
			return fmt.Sprintf("bytes.Equal(%s, %s)", this, that), nil
		}
		p.newFunc(typ)
		return fmt.Sprintf("%s(%s, %s)", funcName(typ, p.qual), this, that), nil
	case *types.Map:
		p.newFunc(typ)
		return fmt.Sprintf("%s(%s, %s)", funcName(typ, p.qual), this, that), nil
	case *types.Named:
		return fmt.Sprintf("%s.Equal(&%s)", this, that), nil
	default: // *Chan, *Tuple, *Signature, *Interface, *types.Basic.Kind() == types.UntypedNil, *Struct
		return "", fmt.Errorf("unsupported type %#v", fieldType)
	}
}

type findFunction struct {
	pkgInfo    *loader.PackageInfo
	funcPrefix string
	typs       []types.Type
}

func (this *findFunction) Visit(node ast.Node) (w ast.Visitor) {
	if call, ok := node.(*ast.CallExpr); ok {
		if fn, ok := call.Fun.(*ast.Ident); ok {
			if strings.HasPrefix(fn.Name, this.funcPrefix) {
				if len(call.Args) != 2 {
					fmt.Fprintf(os.Stderr, "%s does not have two arguments\n", fn.Name)
					return this
				}
				t0 := this.pkgInfo.TypeOf(call.Args[0])
				t1 := this.pkgInfo.TypeOf(call.Args[1])
				if types.Identical(t0, t1) {
					name := strings.TrimPrefix(fn.Name, this.funcPrefix)
					qual := types.RelativeTo(this.pkgInfo.Pkg)
					typeStr := typeName(t0, qual)
					if typeStr != name {
						fmt.Fprintf(os.Stderr, "%s's suffix %s does not match the type %s\n",
							fn.Name, name, typeStr)
						return this
					}
					this.typs = append(this.typs, t0)
					return this
				} else {
					fmt.Fprintf(os.Stderr, "%s has two arguments, but they are of different types %s != %s\n",
						fn.Name, t0, t1)
				}
			}
		}
	}
	return this
}

func findObjectsUsingFunction(pkgInfo *loader.PackageInfo, f *ast.File, funcName string) []types.Type {
	var typs []types.Type
	for _, d := range f.Decls {
		finder := &findFunction{pkgInfo, funcName, nil}
		ast.Walk(finder, d)
		typs = append(typs, finder.typs...)
	}
	return typs
}
