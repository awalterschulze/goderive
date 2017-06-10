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

import "go/types"

type qualifiers struct {
	p        *types.Package
	importer importer
	imported map[*types.Package]Import
}

type importer interface {
	NewImport(path string) Import
}

func newQualifiers(importer importer, p *types.Package) types.Qualifier {
	q := &qualifiers{
		p:        p,
		importer: importer,
		imported: make(map[*types.Package]Import),
	}
	return q.Qualifier
}

func (this *qualifiers) Qualifier(p *types.Package) string {
	if this.p == p {
		return ""
	}
	if _, ok := this.imported[p]; !ok {
		this.imported[p] = this.importer.NewImport(p.Path())
	}
	return this.imported[p]()
}
