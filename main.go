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
	"github.com/awalterschulze/goderive/plugin/compare"
	"github.com/awalterschulze/goderive/plugin/equal"
	"github.com/awalterschulze/goderive/plugin/fmap"
	"github.com/awalterschulze/goderive/plugin/join"
	"github.com/awalterschulze/goderive/plugin/keys"
	"github.com/awalterschulze/goderive/plugin/sorted"
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
		sorted.NewPlugin(),
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
	g, err := derive.NewGenerator(plugins, paths, *autoname, *dedup)
	if err != nil {
		log.Fatal(err)
	}
	if err := g.Generate(); err != nil {
		log.Fatal(err)
	}
}
