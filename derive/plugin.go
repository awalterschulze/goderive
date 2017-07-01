//  Copyright 2017 Walter Schulze
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package derive

import (
	"go/types"
)

type Plugin interface {
	GetPrefix() string
	SetPrefix(string)
	Name() string
	New(typesMap TypesMap, p Printer, deps map[string]Dependency) Generator
}

type Generator interface {
	TypesMap
	Add(name string, typs []types.Type) (string, error)
	Generate(typs []types.Type) error
}

type Dependency interface {
	GetFuncName(typs ...types.Type) string
}

type plugin struct {
	name    string
	prefix  string
	newFunc func(typesMap TypesMap, p Printer, deps map[string]Dependency) Generator
}

func NewPlugin(name, prefix string, newFunc func(typesMap TypesMap, p Printer, deps map[string]Dependency) Generator) Plugin {
	return &plugin{
		name:    name,
		prefix:  prefix,
		newFunc: newFunc,
	}
}

func (g *plugin) New(typesMap TypesMap, p Printer, deps map[string]Dependency) Generator {
	return g.newFunc(typesMap, p, deps)
}

func (g *plugin) GetPrefix() string {
	return g.prefix
}

func (g *plugin) SetPrefix(p string) {
	g.prefix = p
}

func (g *plugin) Name() string {
	return g.name
}
