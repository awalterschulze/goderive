//  Copyright 2021 Jake Son
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
	"strconv"
	"strings"
)

var blackIdentifier = "_"

// RenameBlankIdentifier returns a signature that all black identified are renamed.
func RenameBlankIdentifier(sig *types.Signature, prefixs ...string) *types.Signature {
	prefix := "param_"
	if len(prefixs) > 0 {
		prefix = prefixs[0]
	}
	params := sig.Params()
	if !hasBlankIdentifier(params) {
		return sig
	}
	renamedTuple := rename(params, prefix)
	return types.NewSignature(sig.Recv(), renamedTuple, sig.Results(), sig.Variadic())
}

func hasBlankIdentifier(tup *types.Tuple) bool {
	for i := 0; i < tup.Len(); i++ {
		if tup.At(i).Name() == blackIdentifier {
			return true
		}
	}
	return false
}

func rename(tup *types.Tuple, prefix string) *types.Tuple {
	vars := make([]*types.Var, tup.Len())
	for i := range vars {
		varValue := tup.At(i)
		if varValue.Name() == blackIdentifier || strings.HasPrefix(varValue.Name(), prefix) {
			varValue = types.NewVar(varValue.Pos(), varValue.Pkg(), prefix+strconv.Itoa(i), varValue.Type())
		}
		vars[i] = varValue
	}
	return types.NewTuple(vars...)
}
