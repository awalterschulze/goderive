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

package main

import (
	"fmt"
	"io"
)

type Printer interface {
	P(format string, a ...interface{})
	In()
	Out()
}

type printer struct {
	w      io.Writer
	indent string
}

func newPrinter(w io.Writer) Printer {
	return &printer{w, ""}
}

func (p *printer) P(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(p.w, p.indent+format+"\n", a...); err != nil {
		panic(err)
	}
}

func (p *printer) In() {
	p.indent += "\t"
}

func (p *printer) Out() {
	if len(p.indent) > 0 {
		p.indent = p.indent[1:]
	} else {
		panic("bug in code generator: unindenting more than has been indented")
	}
}
