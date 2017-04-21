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
	p.P(this.sortPkg() + ".Slice(s, func(i, j int) bool { return s[i] < s[j] })")
	p.P("return s")
	p.Out()
	p.P("}")
	return nil
}
