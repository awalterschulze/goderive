package derive

import (
	"go/types"
)

type Generator interface {
	GetPrefix() string
	SetPrefix(string)
	Name() string
	New(typesMap TypesMap, p Printer, deps map[string]Dependency) Plugin
}

type Plugin interface {
	TypesMap
	Add(name string, typs []types.Type) (string, error)
	Generate() error
}

type Dependency interface {
	GetFuncName(typs ...types.Type) string
}

type generator struct {
	name    string
	prefix  string
	newFunc func(typesMap TypesMap, p Printer, deps map[string]Dependency) Plugin
}

func NewGenerator(name, prefix string, newFunc func(typesMap TypesMap, p Printer, deps map[string]Dependency) Plugin) Generator {
	return &generator{
		name:    name,
		prefix:  prefix,
		newFunc: newFunc,
	}
}

func (g *generator) New(typesMap TypesMap, p Printer, deps map[string]Dependency) Plugin {
	return g.newFunc(typesMap, p, deps)
}

func (g *generator) GetPrefix() string {
	return g.prefix
}

func (g *generator) SetPrefix(p string) {
	g.prefix = p
}

func (g *generator) Name() string {
	return g.name
}
