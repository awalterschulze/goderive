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

package test

import (
	"fmt"
	"testing"
)

func TestToError(t *testing.T) {
	m := map[int]string{0: "0"}
	expectKey := func(i int) (a string, b bool) {
		a, b = m[i]
		return
	}
	eFalse := fmt.Errorf("eFalse")
	transformed := deriveToError(eFalse, expectKey)
	str, e := transformed(0)
	if !(e == nil && str == "0") {
		t.Fatalf("expected key 0 %s", e.Error())
	}

	str, e = transformed(1)
	if !(e != nil && e.Error() == eFalse.Error()) {
		t.Fatalf("unexpected key 1 %v", e)
	}
}
