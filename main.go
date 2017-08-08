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
	"flag"
	"log"

	"github.com/awalterschulze/goderive/derive"
	"github.com/awalterschulze/goderive/plugin/all"
	"github.com/awalterschulze/goderive/plugin/any"
	"github.com/awalterschulze/goderive/plugin/bind"
	"github.com/awalterschulze/goderive/plugin/compare"
	"github.com/awalterschulze/goderive/plugin/contains"
	"github.com/awalterschulze/goderive/plugin/copyto"
	"github.com/awalterschulze/goderive/plugin/curry"
	"github.com/awalterschulze/goderive/plugin/equal"
	"github.com/awalterschulze/goderive/plugin/filter"
	"github.com/awalterschulze/goderive/plugin/flip"
	"github.com/awalterschulze/goderive/plugin/fmap"
	"github.com/awalterschulze/goderive/plugin/gostring"
	"github.com/awalterschulze/goderive/plugin/intersect"
	"github.com/awalterschulze/goderive/plugin/join"
	"github.com/awalterschulze/goderive/plugin/keys"
	"github.com/awalterschulze/goderive/plugin/max"
	"github.com/awalterschulze/goderive/plugin/min"
	"github.com/awalterschulze/goderive/plugin/set"
	"github.com/awalterschulze/goderive/plugin/sort"
	"github.com/awalterschulze/goderive/plugin/takewhile"
	"github.com/awalterschulze/goderive/plugin/tuple"
	"github.com/awalterschulze/goderive/plugin/uncurry"
	"github.com/awalterschulze/goderive/plugin/union"
	"github.com/awalterschulze/goderive/plugin/unique"
)

var autoname = flag.Bool("autoname", false, "rename functions that are conflicting with other functions")
var dedup = flag.Bool("dedup", false, "rename functions to functions that are duplicates")

func main() {
	plugins := []derive.Plugin{
		equal.NewPlugin(),
		compare.NewPlugin(),
		fmap.NewPlugin(),
		join.NewPlugin(),
		keys.NewPlugin(),
		sort.NewPlugin(),
		copyto.NewPlugin(),
		set.NewPlugin(),
		min.NewPlugin(),
		max.NewPlugin(),
		contains.NewPlugin(),
		intersect.NewPlugin(),
		union.NewPlugin(),
		filter.NewPlugin(),
		takewhile.NewPlugin(),
		unique.NewPlugin(),
		flip.NewPlugin(),
		curry.NewPlugin(),
		uncurry.NewPlugin(),
		all.NewPlugin(),
		any.NewPlugin(),
		tuple.NewPlugin(),
		bind.NewPlugin(),
		gostring.NewPlugin(),
	}
	flags := make(map[string]*string)
	for _, p := range plugins {
		flags[p.Name()] = flag.String(p.Name()+".prefix", p.GetPrefix(), "set the prefix for "+p.Name()+" functions that should be derived.")
	}
	log.SetFlags(0)
	flag.Parse()
	for _, p := range plugins {
		p.SetPrefix(*(flags[p.Name()]))
	}
	paths := derive.ImportPaths(flag.Args())
	g, err := derive.NewPlugins(plugins, *autoname, *dedup).Load(paths)
	if err != nil {
		log.Fatal(err)
	}
	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
