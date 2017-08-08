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
	"go/parser"
	"go/token"
	"reflect"
	"strings"
	"testing"
)

type gostringer interface {
	GoString() string
}

func TestGoString(t *testing.T) {
	structs := []gostringer{
		&Empty{},
		&BuiltInTypes{},
	}
	for _, this := range structs {
		desc := reflect.TypeOf(this).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			for i := 0; i < 100; i++ {
				this = random(this).(gostringer)
				s := this.GoString()
				content := `package main
				func main() {
				` + s + `
				}
				`
				fset := token.NewFileSet()
				if _, err := parser.ParseFile(fset, "main.go", content, parser.AllErrors); err != nil {
					t.Fatalf("parse error: %v, given input <%s>", err, s)
				}
				if strings.Contains(s, "x0") {
					t.Fatalf("printed a pointer instead of a value in %s", s)
				}
				if i == 0 {
					t.Log(s)
				}
			}
		})
	}
}
