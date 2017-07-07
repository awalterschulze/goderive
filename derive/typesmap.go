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
	"fmt"
	"go/types"
	"strconv"
)

type TypesMap interface {
	SetFuncName(name string, typs ...types.Type) (newName string, err error)
	GetFuncName(typs ...types.Type) string
	Generating(typs ...types.Type)
	ToGenerate() [][]types.Type
	Prefix() string
	TypeString(typ types.Type) string
	IsExternal(typ *types.Named) bool
	Done() bool
}

type typesMap struct {
	qual       types.Qualifier
	prefix     string
	generated  map[string]bool
	funcToTyps map[string][]types.Type
	typss      [][]types.Type
	autoname   bool
	dedup      bool
}

func newTypesMap(qual types.Qualifier, prefix string, autoname bool, dedup bool) TypesMap {
	return &typesMap{
		qual:       qual,
		prefix:     prefix,
		generated:  make(map[string]bool),
		funcToTyps: make(map[string][]types.Type),
		typss:      nil,
		autoname:   autoname,
		dedup:      dedup,
	}
}

func (this *typesMap) Prefix() string {
	return this.prefix
}

func (this *typesMap) TypeString(typ types.Type) string {
	return types.TypeString(types.Default(typ), this.qual)
}

func (this *typesMap) IsExternal(typ *types.Named) bool {
	q := this.qual(typ.Obj().Pkg())
	return q != ""
}

func (this *typesMap) SetFuncName(funcName string, typs ...types.Type) (string, error) {
	if fName, ok := this.nameOf(typs); ok {
		if fName == funcName {
			return funcName, nil
		}
		if this.dedup {
			return fName, nil
		}
		return "", fmt.Errorf("ambigious function names for type %s = (%s | %s)", typs, fName, funcName)
	}
	if ts, ok := this.funcToTyps[funcName]; ok {
		if eq(ts, typs) {
			return funcName, nil
		}
		if this.autoname {
			return this.GetFuncName(typs...), nil
		}
		return "", fmt.Errorf("conflicting function names %s", funcName)
	}
	this.funcToTyps[funcName] = typs
	this.typss = append(this.typss, typs)
	return funcName, nil
}

func (this *typesMap) GetFuncName(typs ...types.Type) string {
	name, ok := this.nameOf(typs)
	if !ok {
		name = this.newName(typs)
		this.SetFuncName(name, typs...)
	}
	return name
}

func (this *typesMap) newName(typs []types.Type) string {
	funcName := this.prefix
	_, exists := this.funcToTyps[funcName]
	i := 0
	name := ""
	if len(typs) > 0 {
		switch t := typs[0].(type) {
		case *types.Named:
			name = t.Obj().Name()
		case *types.Basic:
			switch t.Kind() {
			case types.Bool, types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
				types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
				types.Float32, types.Float64, types.String:
				name = this.TypeString(t)
			}
		}
	}
	for exists {
		if i > len(name) {
			funcName = this.prefix + "_" + name + strconv.Itoa(i)
		} else {
			funcName = this.prefix + "_" + name[:i]
		}
		i++
		_, exists = this.funcToTyps[funcName]
	}
	return funcName
}

func eq(this, that []types.Type) bool {
	if len(this) != len(that) {
		return false
	}
	for i, t := range this {
		if !types.Identical(t, that[i]) {
			return false
		}
	}
	return true
}

func (this *typesMap) nameOf(typs []types.Type) (string, bool) {
	for _, t := range typs {
		if n, ok := t.(*types.Named); ok {
			this.qual(n.Obj().Pkg())
		}
	}
	for name, ts := range this.funcToTyps {
		if eq(typs, ts) {
			return name, true
		}
	}
	return "", false
}

func (this *typesMap) Generating(typs ...types.Type) {
	name, ok := this.nameOf(typs)
	if !ok {
		panic("wtf")
	}
	this.generated[name] = true
}

func (this *typesMap) isGenerated(typs []types.Type) bool {
	name, ok := this.nameOf(typs)
	if !ok {
		return false
	}
	return this.generated[name]
}

func (this *typesMap) ToGenerate() [][]types.Type {
	typss := make([][]types.Type, 0, len(this.typss))
	for i, typs := range this.typss {
		if !this.isGenerated(typs) {
			typss = append(typss, this.typss[i])
		}
	}
	return typss
}

func (this *typesMap) Done() bool {
	for _, typs := range this.typss {
		if !this.isGenerated(typs) {
			return false
		}
	}
	return true
}
