package sort

import (
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

func NewPlugin() derive.Plugin {
	return derive.NewPlugin("sort", "deriveSort", New)
}

func New(typesMap derive.TypesMap, p derive.Printer, deps map[string]derive.Dependency) derive.Generator {
	return &gen{
		TypesMap: typesMap,
		printer:  p,
		sortPkg:  p.NewImport("sort"),
		compare:  deps["compare"],
	}
}

type gen struct {
	derive.TypesMap
	printer derive.Printer
	sortPkg derive.Import
	compare derive.Dependency
}

func (this *gen) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
}

func (this *gen) Generate() error {
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

func (this *gen) genFuncFor(typ *types.Slice) error {
	p := this.printer
	this.Generating(typ)
	typeStr := this.TypeString(typ)
	p.P("")
	p.P("func %s(list %s) %s {", this.GetFuncName(typ), typeStr, typeStr)
	p.In()
	etyp := typ.Elem()
	switch ttyp := etyp.Underlying().(type) {
	case *types.Basic:
		switch ttyp.Kind() {
		case types.String:
			p.P(this.sortPkg() + ".Strings(list)")
		case types.Float64:
			p.P(this.sortPkg() + ".Float64s(list)")
		case types.Int:
			p.P(this.sortPkg() + ".Ints(list)")
		case types.Complex64, types.Complex128, types.Bool:
			p.P(this.sortPkg() + ".Slice(list, func(i, j int) bool { return " + this.compare.GetFuncName(ttyp) + "(list[i], list[j]) < 0 })")
		default:
			p.P(this.sortPkg() + ".Slice(list, func(i, j int) bool { return list[i] < list[j] })")
		}
	case *types.Pointer, *types.Struct, *types.Slice, *types.Array, *types.Map:
		p.P(this.sortPkg() + ".Slice(list, func(i, j int) bool { return " + this.compare.GetFuncName(etyp) + "(list[i], list[j]) < 0 })")
	default:
		return fmt.Errorf("unsupported compare type: %s", this.TypeString(typ))
	}
	p.P("return list")
	p.Out()
	p.P("}")
	return nil
}
