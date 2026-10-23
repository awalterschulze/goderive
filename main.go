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

	"go.awalterschulze.org/goderive/derive"
	"go.awalterschulze.org/goderive/plugin/all"
	"go.awalterschulze.org/goderive/plugin/any"
	"go.awalterschulze.org/goderive/plugin/apply"
	"go.awalterschulze.org/goderive/plugin/clone"
	"go.awalterschulze.org/goderive/plugin/compare"
	"go.awalterschulze.org/goderive/plugin/compose"
	"go.awalterschulze.org/goderive/plugin/contains"
	"go.awalterschulze.org/goderive/plugin/curry"
	"go.awalterschulze.org/goderive/plugin/deepcopy"
	"go.awalterschulze.org/goderive/plugin/do"
	"go.awalterschulze.org/goderive/plugin/dup"
	"go.awalterschulze.org/goderive/plugin/equal"
	"go.awalterschulze.org/goderive/plugin/filter"
	"go.awalterschulze.org/goderive/plugin/flip"
	"go.awalterschulze.org/goderive/plugin/fmap"
	"go.awalterschulze.org/goderive/plugin/gostring"
	"go.awalterschulze.org/goderive/plugin/hash"
	"go.awalterschulze.org/goderive/plugin/intersect"
	"go.awalterschulze.org/goderive/plugin/join"
	"go.awalterschulze.org/goderive/plugin/keys"
	"go.awalterschulze.org/goderive/plugin/max"
	"go.awalterschulze.org/goderive/plugin/mem"
	"go.awalterschulze.org/goderive/plugin/min"
	"go.awalterschulze.org/goderive/plugin/pipeline"
	"go.awalterschulze.org/goderive/plugin/set"
	"go.awalterschulze.org/goderive/plugin/sort"
	"go.awalterschulze.org/goderive/plugin/takewhile"
	"go.awalterschulze.org/goderive/plugin/toerror"
	"go.awalterschulze.org/goderive/plugin/traverse"
	"go.awalterschulze.org/goderive/plugin/tuple"
	"go.awalterschulze.org/goderive/plugin/uncurry"
	"go.awalterschulze.org/goderive/plugin/union"
	"go.awalterschulze.org/goderive/plugin/unique"
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
		toerror.NewPlugin(),
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
		hash.NewPlugin(),
		mem.NewPlugin(),
		traverse.NewPlugin(),
		apply.NewPlugin(),
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
