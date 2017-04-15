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
	printer Printer
	sortPkg Import
	compare TypesMap
}

func newSortedKeys(p Printer, qual types.Qualifier, typesMap TypesMap, compare TypesMap) (*sortedKeys, error) {
	return &sortedKeys{
		TypesMap: typesMap,
		qual:     qual,
		printer:  p,
		sortPkg:  p.NewImport("sort"),
		compare:  compare,
	}, nil
}

func (this *sortedKeys) Generate(pkgInfo *loader.PackageInfo, prefix string, call *ast.CallExpr) (bool, error) {
	fn, ok := call.Fun.(*ast.Ident)
	if !ok {
		return false, nil
	}
	if !strings.HasPrefix(fn.Name, prefix) {
		return false, nil
	}
	if len(call.Args) != 1 {
		return false, fmt.Errorf("%s does not have one argument", fn.Name)
	}
	typ := pkgInfo.TypeOf(call.Args[0])
	if err := this.SetFuncName(typ, fn.Name); err != nil {
		return false, err
	}

	for _, typ := range this.ToGenerate() {
		mapType, ok := typ.(*types.Map)
		if !ok {
			return false, fmt.Errorf("%s, an argument to %s, is not of type map", this.GetFuncName(typ), typ)
		}
		if err := this.genFuncFor(mapType); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (this *sortedKeys) genFuncFor(typ *types.Map) error {
	p := this.printer
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
