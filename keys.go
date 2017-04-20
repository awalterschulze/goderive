package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

var keysPrefix = flag.String("keys.prefix", "deriveKeys", "set the prefix for keys functions that should be derived.")

type keys struct {
	TypesMap
	qual    types.Qualifier
	prefix  string
	printer Printer
}

func newKeys(typesMap TypesMap, qual types.Qualifier, prefix string, p Printer) *keys {
	return &keys{
		TypesMap: typesMap,
		qual:     qual,
		prefix:   prefix,
		printer:  p,
	}
}

func (this *keys) Add(pkgInfo *loader.PackageInfo, call *ast.CallExpr) (bool, error) {
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

func (this *keys) Generate() error {
	for _, typs := range this.ToGenerate() {
		typ := typs[0]
		mapType, ok := typ.(*types.Map)
		if !ok {
			return fmt.Errorf("%s, the first argument, %s, is not of type map", this.GetFuncName(typ), typ)
		}
		if err := this.genFuncFor(mapType); err != nil {
			return err
		}
	}
	return nil
}

func (this *keys) genFuncFor(typ *types.Map) error {
	p := this.printer
	this.Generating(typ)
	typeStr := types.TypeString(typ, this.qual)
	keyType := typ.Key()
	keyTypeStr := types.TypeString(keyType, this.qual)
	p.P("")
	p.P("func %s(m %s) []%s {", this.GetFuncName(typ), typeStr, keyTypeStr)
	p.In()
	p.P("keys := make([]%s, 0, len(m))", keyTypeStr)
	p.P("for key, _ := range m {")
	p.In()
	p.P("keys = append(keys, key)")
	p.Out()
	p.P("}")
	p.P("return keys")
	p.Out()
	p.P("}")
	return nil
}
