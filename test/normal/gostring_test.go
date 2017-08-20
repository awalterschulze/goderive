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
	"bytes"
	"encoding/gob"
	"fmt"
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
		&SliceOfBuiltInTypes{},
		&SliceOfPtrToBuiltInTypes{},
		&ArrayOfBuiltInTypes{},
		&ArrayOfPtrToBuiltInTypes{},
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
	f.WriteString("\t\"encoding/gob\"\n")
	f.WriteString("\t\"bytes\"\n")
	f.WriteString("\t\"reflect\"\n")
	f.WriteString(")\n")
	f.WriteString("\n")
	f.WriteString("func TestGeneratedGoString(t *testing.T) {\n")
	for _, empty := range structs {
		desc := reflect.TypeOf(empty).Elem().Name()
		t.Run(desc, func(t *testing.T) {
			first := true
			for i := 0; i < 100; i++ {
				this := random(empty).(gostringer)
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
				if ss, ok := this.(*SliceOfPtrToBuiltInTypes); i == 0 && ok {
					t.Logf("%#v", *ss)
					t.Log(s)
				}
				if first {
					// If gob has not been able to encode any of 10 random variables,
					// then I guess its time to move on to a simpler test.
					if i == 10 {
						first = false
						fmt.Fprintf(f, "t.Run(%q, func(t *testing.T) {\n", desc)
						f.WriteString(s)
						f.WriteString("})\n")
					} else if !reflect.ValueOf(this).IsNil() {
						buf := bytes.NewBuffer(nil)
						enc := gob.NewEncoder(buf)
						// Suprisingly many things that gob cannot encode.
						if err := enc.Encode(this); err == nil {
							first = false
							fmt.Fprintf(f, "t.Run(%q, func(t *testing.T) {\n", desc)

							fmt.Fprintf(f, "data := %#v\n", buf.Bytes())

							f.WriteString("gotbeforegob := " + s)
							f.WriteString("buf := bytes.NewBuffer(data)\n")

							f.WriteString("dec := gob.NewDecoder(buf)\n")
							fmt.Fprintf(f, "want := %#v\n", empty)
							f.WriteString("if err := dec.Decode(want); err != nil {\n")
							f.WriteString("t.Fatal(err)\n")
							f.WriteString("}\n")

							//gob sees nil and empty slices as the same thing.
							f.WriteString("equalizer := bytes.NewBuffer(nil)\n")
							f.WriteString("enc := gob.NewEncoder(equalizer)\n")
							f.WriteString("dec = gob.NewDecoder(equalizer)\n")
							f.WriteString("enc.Encode(gotbeforegob)\n")
							fmt.Fprintf(f, "got := %#v\n", empty)
							f.WriteString("if err := dec.Decode(got); err != nil {\n")
							f.WriteString("t.Fatal(err)\n")
							f.WriteString("}\n")

							f.WriteString("if !reflect.DeepEqual(got, want) {\n")
							f.WriteString("if got != nil && want != nil {\n")
							f.WriteString("t.Fatalf(\"got %#v != want %#v\", *got, *want)\n")
							f.WriteString("} else {\n")
							f.WriteString("t.Fatalf(\"got %#v != want %#v\", got, want)\n")
							f.WriteString("}\n")
							f.WriteString("}\n")
							f.WriteString("})\n")
						}
					}
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
