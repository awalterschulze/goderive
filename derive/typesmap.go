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
	reserved   map[string]struct{}
	autoname   bool
	dedup      bool
}

func newTypesMap(qual types.Qualifier, prefix string, reserved map[string]struct{}, autoname bool, dedup bool) TypesMap {
	return &typesMap{
		qual:       qual,
		prefix:     prefix,
		generated:  make(map[string]bool),
		funcToTyps: make(map[string][]types.Type),
		typss:      nil,
		reserved:   reserved,
		autoname:   autoname,
		dedup:      dedup,
	}
}

func (tm *typesMap) Prefix() string {
	return tm.prefix
}

func (tm *typesMap) TypeString(typ types.Type) string {
	return types.TypeString(types.Default(typ), tm.qual)
}

func (tm *typesMap) TypeStringBypass(typ types.Type) string {
	return types.TypeString(types.Default(typ), bypassQual)
}

func (tm *typesMap) IsExternal(typ *types.Named) bool {
	q := tm.qual(typ.Obj().Pkg())
	return q != ""
}

func (tm *typesMap) SetFuncName(funcName string, typs ...types.Type) (string, error) {
	// log.Printf("SetFuncName: %s(%v)", funcName, typs)
	if fName, ok := tm.nameOf(typs); ok {
		if fName == funcName {
			return funcName, nil
		}
		if tm.dedup {
			return fName, nil
		}
		return "", fmt.Errorf("ambigious function names for type %s = (%s | %s)", typs, fName, funcName)
	}
	if ts, ok := tm.funcToTyps[funcName]; ok {
		if eq(ts, typs) {
			return funcName, nil
		}
		if tm.autoname {
			return tm.GetFuncName(typs...), nil
		}
		return "", fmt.Errorf("conflicting function names %s(%v) and %s(%v)", funcName, ts, funcName, typs)
	}
	tm.funcToTyps[funcName] = typs
	tm.typss = append(tm.typss, typs)
	return funcName, nil
}

func (tm *typesMap) GetFuncName(typs ...types.Type) string {
	// log.Printf("GetFuncName: %v", typs)
	name, ok := tm.nameOf(typs)
	if !ok {
		name = tm.newName(typs)
		tm.SetFuncName(name, typs...)
	}
	// log.Printf("GotFuncName: %s(%v)", name, typs)
	return name
}

func (tm *typesMap) newName(typs []types.Type) string {
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
				name = tm.TypeString(t)
			}
		}
	}
	i := 0
	funcName := tm.prefix
	_, exists := tm.funcToTyps[funcName]
	_, isreserved := tm.reserved[funcName]
	for exists || isreserved {
		if i > len(name) {
			funcName = tm.prefix + "_" + name + strconv.Itoa(i)
		} else {
			funcName = tm.prefix + "_" + name[:i]
		}
		i++
		_, exists = tm.funcToTyps[funcName]
		_, isreserved = tm.reserved[funcName]
	}
	return funcName
}

func eq(this, that []types.Type) bool {
	if len(this) != len(that) {
		return false
	}
	for i, t := range this {
		if !types.AssignableTo(types.Default(t), types.Default(that[i])) {
			return false
		}
	}
	return true
}

func (tm *typesMap) nameOf(typs []types.Type) (string, bool) {
	for _, t := range typs {
		if n, ok := t.(*types.Named); ok {
			pkg := n.Obj().Pkg()
			if pkg != nil {
				tm.qual(pkg)
			}
		}
	}
	for name, ts := range tm.funcToTyps {
		if eq(typs, ts) {
			return name, true
		}
	}
	return "", false
}

func (tm *typesMap) Generating(typs ...types.Type) {
	name, ok := tm.nameOf(typs)
	if !ok {
		panic(fmt.Sprintf("generating unknown %s for types: %v", tm.prefix, typs))
	}
	tm.generated[name] = true
}

func (tm *typesMap) isGenerated(typs []types.Type) bool {
	name, ok := tm.nameOf(typs)
	if !ok {
		return false
	}
	return tm.generated[name]
}

func (tm *typesMap) ToGenerate() [][]types.Type {
	typss := make([][]types.Type, 0, len(tm.typss))
	for i, typs := range tm.typss {
		if !tm.isGenerated(typs) {
			typss = append(typss, tm.typss[i])
		}
	}
	return typss
}

func (tm *typesMap) Done() bool {
	for _, typs := range tm.typss {
		if !tm.isGenerated(typs) {
			return false
		}
	}
	return true
}
