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
	"os"
	"os/exec"
	"reflect"
	"testing"
)

type gostringer interface {
	GoString() string
}

func TestGoString(t *testing.T) {
	structs := []gostringer{
		&Empty{},
		&BuiltInTypes{},
		&PtrToBuiltInTypes{},
	}
	filename := "gostring_gen_test.go"
	f, err := os.Create(filename)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("package test\n")
	f.WriteString("\n")
	f.WriteString("import (\n")
	f.WriteString("\t\"testing\"\n")
	f.WriteString(")\n")
	f.WriteString("\n")
	f.WriteString("func TestGeneratedGoString(t *testing.T) {\n")
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
				if i == 0 {
					f.WriteString(s)
				}
			}
		})
	}
	f.WriteString("}\n")
	f.Close()
	gofmtcmd := exec.Command("gofmt", "-l", "-s", "-w", filename)
	if o, err := gofmtcmd.CombinedOutput(); err != nil {
		t.Fatalf("%q, error: %v", o, err)
	}
	testcmd := exec.Command("go", "test", "-v", "-run", "TestGeneratedGoString")
	if o, err := testcmd.CombinedOutput(); err != nil {
		t.Fatalf("%s, error: %v", o, err)
	}
}
