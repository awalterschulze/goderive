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

func generateSortedKeys(p Printer, pkgInfo *loader.PackageInfo, prefix string, strict bool, calls []*ast.CallExpr) error {
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
			return fmt.Errorf("%s does not have one argument\n", fn.Name)
		}
		typ := pkgInfo.TypeOf(call.Args[0])
		if err := typesMap.SetFuncName(typ, fn.Name); err != nil {
			return err
		}
	}

	sortedKeys := newSortedKeys(p, typesMap, qual)

	for _, typ := range typesMap.List() {
		mapType, ok := typ.(*types.Map)
		if !ok {
			return fmt.Errorf("%s, an argument to %s, is not of type map", typesMap.GetFuncName(typ), typ)
		}
		if err := sortedKeys.genFuncFor(mapType); err != nil {
			return err
		}
	}
	return nil
}

type sortedKeys struct {
	printer  Printer
	typesMap TypesMap
	qual     types.Qualifier
	sortPkg  Import
}

func newSortedKeys(printer Printer, typesMap TypesMap, qual types.Qualifier) *sortedKeys {
	return &sortedKeys{
		printer:  printer,
		typesMap: typesMap,
		qual:     qual,
		sortPkg:  printer.NewImport("sort"),
	}
}

func (this *sortedKeys) genFuncFor(typ *types.Map) error {
	p := this.printer
	this.typesMap.Generating(typ)
	typeStr := types.TypeString(typ, this.qual)
	keyType := typ.Key()
	keyTypeStr := types.TypeString(keyType, this.qual)
	p.P("")
	p.P("func %s(m %s) []%s {", this.typesMap.GetFuncName(typ), typeStr, keyTypeStr)
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
