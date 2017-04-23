package main

import (
	"flag"
	"fmt"
	"go/types"
)

var sortedPrefix = flag.String("sorted.prefix", "deriveSorted", "set the prefix for sorted functions that should be derived.")

type sorted struct {
	TypesMap
	printer Printer
	sortPkg Import
	compare Plugin
}

func newSorted(typesMap TypesMap, p Printer, compareTypesMap Plugin) *sorted {
	return &sorted{
		TypesMap: typesMap,
		printer:  p,
		sortPkg:  p.NewImport("sort"),
		compare:  compareTypesMap,
	}
}

func (this *sorted) Name() string {
	return "sorted"
}

func (this *sorted) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *sorted) Generate() error {
	for _, typs := range this.ToGenerate() {
		typ := typs[0]
		sliceType, ok := typ.(*types.Slice)
		if !ok {
			return fmt.Errorf("%s, the first argument, %s, is not of type slice", this.GetFuncName(typ), this.TypeString(typ))
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
	typeStr := this.TypeString(typ)
	p.P("")
	p.P("func %s(s %s) %s {", this.GetFuncName(typ), typeStr, typeStr)
	p.In()
	etyp := typ.Elem()
	switch ttyp := etyp.Underlying().(type) {
	case *types.Basic:
		switch ttyp.Kind() {
		case types.String:
			p.P(this.sortPkg() + ".Strings(s)")
		case types.Float64:
			p.P(this.sortPkg() + ".Float64s(s)")
		case types.Int:
			p.P(this.sortPkg() + ".Ints(s)")
		case types.Complex64, types.Complex128, types.Bool:
			p.P(this.sortPkg() + ".Slice(s, func(i, j int) bool { return " + this.compare.GetFuncName(ttyp) + "(s[i], s[j]) < 0 })")
		default:
			p.P(this.sortPkg() + ".Slice(s, func(i, j int) bool { return s[i] < s[j] })")
		}
	case *types.Pointer, *types.Struct, *types.Slice, *types.Array, *types.Map:
		p.P(this.sortPkg() + ".Slice(s, func(i, j int) bool { return " + this.compare.GetFuncName(etyp) + "(s[i], s[j]) < 0 })")
	default:
		return fmt.Errorf("unsupported compare type: %#v", typ)
	}
	p.P("return s")
	p.Out()
	p.P("}")
	return nil
}
