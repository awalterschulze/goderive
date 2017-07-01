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
	"strings"
)

type TypesMap interface {
	SetFuncName(name string, typs ...types.Type) (newName string, err error)
	GetFuncName(typs ...types.Type) string
	Generating(typs ...types.Type)
	ToGenerate() [][]types.Type
	Prefix() string
	TypeString(typ types.Type) string
	Done() bool
}

type typesMap struct {
	qual       types.Qualifier
	prefix     string
	generated  map[string]bool
	funcToTyps map[string]string
	typsToFunc map[string]string
	typss      [][]types.Type
	autoname   bool
	dedup      bool
}

func newTypesMap(qual types.Qualifier, prefix string, autoname bool, dedup bool) TypesMap {
	return &typesMap{
		qual:       qual,
		prefix:     prefix,
		generated:  make(map[string]bool),
		funcToTyps: make(map[string]string),
		typsToFunc: make(map[string]string),
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

func (this *typesMap) SetFuncName(funcName string, typs ...types.Type) (string, error) {
	typsName := nameOf(typs, this.qual)
	if fName, ok := this.typsToFunc[typsName]; ok {
		if fName == funcName {
			return funcName, nil
		}
		if this.dedup {
			return fName, nil
		}
		return "", fmt.Errorf("ambigious function names for type %s = (%s | %s)", typs, fName, funcName)
	}
	if tName, ok := this.funcToTyps[funcName]; ok {
		if tName == typsName {
			return funcName, nil
		}
		if this.autoname {
			return this.GetFuncName(typs...), nil
		}
		return "", fmt.Errorf("conflicting function names %s = (%s | %s)", funcName, tName, typsName)
	}
	if _, ok := this.generated[typsName]; !ok {
		this.generated[typsName] = false
	}
	this.typsToFunc[typsName] = funcName
	this.funcToTyps[funcName] = typsName
	this.typss = append(this.typss, typs)
	return funcName, nil
}

func (this *typesMap) GetFuncName(typs ...types.Type) string {
	name := nameOf(typs, this.qual)
	if f, ok := this.typsToFunc[name]; ok {
		return f
	}
	funcName := this.funcOf(typs)
	_, exists := this.funcToTyps[funcName]
	for exists {
		funcName += "_"
		_, exists = this.funcToTyps[funcName]
	}
	this.SetFuncName(funcName, typs...)
	return funcName
}

func (this *typesMap) Generating(typs ...types.Type) {
	name := nameOf(typs, this.qual)
	this.generated[name] = true
}

func (this *typesMap) isGenerated(typs []types.Type) bool {
	name := nameOf(typs, this.qual)
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

func nameOf(typs []types.Type, qual types.Qualifier) string {
	ss := make([]string, len(typs))
	for i, typ := range typs {
		ss[i] = typeName(types.Default(typ), qual)
	}
	return strings.Join(ss, ",")
}

func (this *typesMap) funcOf(typs []types.Type) string {
	return this.prefix + strings.Replace(nameOf(typs, this.qual), "$", "", -1)
}

func typeName(typ types.Type, qual types.Qualifier) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return "PtrTo" + typeName(t.Elem(), qual)
	case *types.Array:
		sizeStr := strconv.Itoa(int(t.Len()))
		return "Array" + sizeStr + "Of" + typeName(t.Elem(), qual)
	case *types.Slice:
		return "SliceOf" + typeName(t.Elem(), qual)
	case *types.Map:
		return "MapOf" + typeName(t.Key(), qual) + "To" + typeName(t.Elem(), qual)
	case *types.Signature:
		params := make([]types.Type, t.Params().Len())
		for i := range params {
			params[i] = t.Params().At(i).Type()
		}
		returns := make([]types.Type, t.Results().Len())
		for i := range returns {
			returns[i] = t.Results().At(i).Type()
		}
		return "FuncOf" + nameOf(params, qual) + "___" + nameOf(returns, qual)
	}
	// The dollar helps to make sure that typenames cannot be faked by the user.
	return "$" + strings.Replace(types.TypeString(typ, qual), ".", "_", -1)
}
