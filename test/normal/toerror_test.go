//  Copyright 2019 Ingun Jon
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
	"fmt"
	"testing"
	"time"
)

type LocalType struct{}

func true0() bool {
	return true
}
func false0() bool {
	return false
}
func true1() (int, bool) {
	return 0, true
}
func true2(a int) (int, bool) {
	return a, true
}
func true3(a, b int) (int, bool) {
	return a + b, true
}
func true4(a, b int) (int, int, bool) {
	return a, b, true
}
func true5(lt *LocalType) (*LocalType, bool) {
	return lt, true
}
func true6(t *time.Time) (*time.Time, bool) {
	return t, true
}

// func fail
func TestToError(t *testing.T) {
	e := fmt.Errorf("error")
	if r := deriveToError0(e, true0)(); !(r == nil) {
		t.Fatal()
	}
	if r := deriveToError0(e, false0)(); !(r == e) {
		t.Fatal()
	}
	if r0, r1 := deriveToError1(e, true1)(); !(r0 == 0 && r1 == nil) {
		t.Fatal()
	}
	if r0, r1 := deriveToError2(e, true2)(1); !(r0 == 1 && r1 == nil) {
		t.Fatal()
	}
	if r0, r1 := deriveToError3(e, true3)(1, 2); !(r0 == 3 && r1 == nil) {
		t.Fatal()
	}
	if r0, r1, r2 := deriveToError4(e, true4)(1, 2); !(r0 == 1 && r1 == 2 && r2 == nil) {
		t.Fatal()
	}
	lt := LocalType{}
	if r0, r1 := deriveToError5(e, true5)(&lt); !(r0 == &lt && r1 == nil) {
		t.Fatal()
	}
	tm := time.Now()
	if r0, r1 := deriveToError6(e, true6)(&tm); !(r0 == &tm && r1 == nil) {
		t.Fatal()
	}
}
