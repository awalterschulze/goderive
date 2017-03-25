package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

var sortedKeysPrefix = flag.String("sortedkeys.prefix", "deriveSortedKeys", "set the prefix for sorted keys functions that should be derived.")

type sortedKeys struct {
	TypesMap
	qual    types.Qualifier
	sortPkg Import
}

func newSortedKeys(pkgInfo *loader.PackageInfo, prefix string, calls []*ast.CallExpr) (*sortedKeys, error) {
	qual := types.RelativeTo(pkgInfo.Pkg)
	typesMap := newTypesMap(qual, prefix)
	for _, call := range calls {
		fn, ok := call.Fun.(*ast.Ident)
		if !ok {
			continue
		}
		if !strings.HasPrefix(fn.Name, prefix) {
			continue
		}
		if len(call.Args) != 1 {
			return nil, fmt.Errorf("%s does not have one argument\n", fn.Name)
		}
		typ := pkgInfo.TypeOf(call.Args[0])
		if err := typesMap.SetFuncName(typ, fn.Name); err != nil {
			return nil, err
		}
	}
	return &sortedKeys{
		TypesMap: typesMap,
		qual:     qual,
	}, nil
}

func (this *sortedKeys) Generate(p Printer) error {
	if this.sortPkg == nil {
		this.sortPkg = p.NewImport("sort")
	}
	for _, typ := range this.ToGenerate() {
		mapType, ok := typ.(*types.Map)
		if !ok {
			return fmt.Errorf("%s, an argument to %s, is not of type map", this.GetFuncName(typ), typ)
		}
		if err := this.genFuncFor(p, mapType); err != nil {
			return err
		}
	}
	return nil
}

func (this *sortedKeys) genFuncFor(p Printer, typ *types.Map) error {
	this.Generating(typ)
	typeStr := types.TypeString(typ, this.qual)
	keyType := typ.Key()
	keyTypeStr := types.TypeString(keyType, this.qual)
	p.P("")
	p.P("func %s(m %s) []%s {", this.GetFuncName(typ), typeStr, keyTypeStr)
	p.In()
	p.P("var keys []%s", keyTypeStr)
	p.P("for key, _ := range m {")
	p.In()
	p.P("keys = append(keys, key)")
	p.Out()
	p.P("}")
	p.P(this.sortPkg() + ".Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })")
	p.P("return keys")
	p.Out()
	p.P("}")
	return nil
}
