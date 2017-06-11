package keys

import (
	"flag"
	"fmt"
	"go/types"

	"github.com/awalterschulze/goderive/derive"
)

var Prefix = flag.String("keys.prefix", "deriveKeys", "set the prefix for keys functions that should be derived.")

type keys struct {
	derive.TypesMap
	printer derive.Printer
}

func New(typesMap derive.TypesMap, p derive.Printer) *keys {
	return &keys{
		TypesMap: typesMap,
		printer:  p,
	}
}

func (this *keys) Name() string {
	return "keys"
}

func (this *keys) Add(name string, typs []types.Type) (string, error) {
	if len(typs) != 1 {
		return "", fmt.Errorf("%s does not have one argument", name)
	}
	return this.SetFuncName(name, typs[0])
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
	typeStr := this.TypeString(typ)
	keyType := typ.Key()
	keyTypeStr := this.TypeString(keyType)
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
