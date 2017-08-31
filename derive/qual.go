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
	"go/types"
)

type qual struct {
	p        *types.Package
	importer importer
	imported map[*types.Package]Import
}

type importer interface {
	NewImport(name, path string) Import
}

func newQualifier(importer importer, p *types.Package) types.Qualifier {
	q := &qual{
		p:        p,
		importer: importer,
		imported: make(map[*types.Package]Import),
	}
	return q.Qualifier
}

func (q *qual) Qualifier(p *types.Package) string {
	if q.p == p {
		return ""
	}
	if _, ok := q.imported[p]; !ok {
		q.imported[p] = q.importer.NewImport(p.Name(), p.Path())
	}
	return q.imported[p]()
}

func bypassQual(p *types.Package) string {
	return p.Name()
}
