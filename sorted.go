package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

var sortedPrefix = flag.String("sorted.prefix", "deriveSorted", "set the prefix for sorted functions that should be derived.")

type sorted struct {
	TypesMap
	qual    types.Qualifier
	prefix  string
	printer Printer
	sortPkg Import
	compare Plugin
}

func newSorted(typesMap TypesMap, qual types.Qualifier, prefix string, p Printer, compareTypesMap Plugin) *sorted {
	return &sorted{
		TypesMap: typesMap,
		qual:     qual,
		prefix:   prefix,
		printer:  p,
		sortPkg:  p.NewImport("sort"),
		compare:  compareTypesMap,
	}
}

func (this *sorted) Add(pkgInfo *loader.PackageInfo, call *ast.CallExpr) (bool, error) {
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil
	}
	if !strings.HasPrefix(fn.Name, this.prefix) {
		return false, nil
	}
	if len(call.Args) != 1 {
		return false, fmt.Errorf("%s does not have one argument", fn.Name)
	}
	typ := pkgInfo.TypeOf(call.Args[0])
	if typ == nil {
		return false, nil
	}
	if err := this.SetFuncName(fn.Name, typ); err != nil {
		return false, err
	}
	return true, nil
}

func (this *sorted) Generate() error {
	for _, typs := range this.ToGenerate() {
		typ := typs[0]
		sliceType, ok := typ.(*types.Slice)
		if !ok {
			return fmt.Errorf("%s, the first argument, %s, is not of type slice", this.GetFuncName(typ), types.TypeString(typ, this.qual))
		}
		if err := this.genFuncFor(sliceType); err != nil {
			return err
		}
	}
	return nil
}

func (this *sorted) genFuncFor(typ *types.Slice) error {
	p := this.printer
	this.Generating(typ)
	typeStr := types.TypeString(typ, this.qual)
	p.P("")
	p.P("func %s(s %s) %s {", this.GetFuncName(typ), typeStr, typeStr)
	p.In()
	p.P(this.sortPkg() + ".Slice(s, func(i, j int) bool { return s[i] < s[j] })")
	p.P("return s")
	p.Out()
	p.P("}")
	return nil
}
