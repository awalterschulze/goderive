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

// Package main implements the goderive binary.
// This pulls in all the plugins, parses the flags and runs the generators using the derive library.
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/awalterschulze/goderive/derive"
	"github.com/awalterschulze/goderive/plugin/all"
	"github.com/awalterschulze/goderive/plugin/any"
	"github.com/awalterschulze/goderive/plugin/clone"
	"github.com/awalterschulze/goderive/plugin/compare"
	"github.com/awalterschulze/goderive/plugin/compose"
	"github.com/awalterschulze/goderive/plugin/contains"
	"github.com/awalterschulze/goderive/plugin/curry"
	"github.com/awalterschulze/goderive/plugin/deepcopy"
	"github.com/awalterschulze/goderive/plugin/do"
	"github.com/awalterschulze/goderive/plugin/dup"
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
	"github.com/awalterschulze/goderive/plugin/pipeline"
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
var prefix = flag.String("prefix", "derive", "prefix of all functions")
var pluginprefix = flag.String("pluginprefix", "", "used to override function prefixes.  The input is a comma separated list of function are prefix pairs.  For example equal=deriveEqual,copyto=copyTo,fmap=fmap,")

func main() {
	plugins := []derive.Plugin{
		equal.NewPlugin(),
		compare.NewPlugin(),
		fmap.NewPlugin(),
		join.NewPlugin(),
		keys.NewPlugin(),
		sort.NewPlugin(),
		deepcopy.NewPlugin(),
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
		gostring.NewPlugin(),
		compose.NewPlugin(),
		do.NewPlugin(),
		pipeline.NewPlugin(),
		dup.NewPlugin(),
		clone.NewPlugin(),
	}
	log.SetFlags(0)
	flag.Parse()
	overridePrefixes := make(map[string]string)
	if len(*pluginprefix) > 0 {
		pairs := strings.Split(*pluginprefix, ",")
		for _, pair := range pairs {
			ss := strings.Split(pair, "=")
			if len(ss) != 2 {
				log.Fatalf("invalid syntax for plugin prefix <%s>", pair)
			}
			overridePrefixes[ss[0]] = ss[1]
		}
	}
	for _, p := range plugins {
		pluginprefix := p.GetPrefix()
		pluginprefix = strings.Replace(pluginprefix, "derive", *prefix, 1)
		newprefix, override := overridePrefixes[p.Name()]
		if override {
			pluginprefix = newprefix
		}
		p.SetPrefix(pluginprefix)
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
