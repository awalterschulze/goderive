package derive

import (
	"go/types"
)

type Generator interface {
	Prefix() string
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
