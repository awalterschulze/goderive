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

import "go/types"

// IsError returns whether a type implements the Error interface.
func IsError(t types.Type) bool {
	typ, ok := t.(*types.Named)
	if !ok {
		return false
	}
	if typ.Obj().Name() == "error" {
		return true
	}
	for i := 0; i < typ.NumMethods(); i++ {
		meth := typ.Method(i)
		if meth.Name() != "Error" {
			continue
		}
		sig, ok := meth.Type().(*types.Signature)
		if !ok {
			// impossible, but lets check anyway
			continue
		}
		if sig.Params().Len() != 0 {
			continue
		}
		res := sig.Results()
		if res.Len() != 1 {
			continue
		}
		b, ok := res.At(0).Type().(*types.Basic)
		if !ok {
			continue
		}
		if b.Kind() != types.String {
			continue
		}
		return true
	}
	return false
}

// Zero returns the zero value as a string, for a given type.
func Zero(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Basic:
		switch t.Kind() {
		case types.String:
			return `""`
		default:
			return "0"
		}
	}
	return "nil"
}

func IsComparable(tt types.Type) bool {
	t := tt.Underlying()
	switch typ := t.(type) {
	case *types.Basic:
		return typ.Kind() != types.UntypedNil
	case *types.Struct:
		for i := 0; i < typ.NumFields(); i++ {
			f := typ.Field(i)
			ft := f.Type()
			if !IsComparable(ft) {
				return false
			}
		}
		return true
	case *types.Array:
		return IsComparable(typ.Elem())
	}
	return false
}
