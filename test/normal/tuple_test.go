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
	"errors"
	"io/ioutil"
	"testing"
)

func TestTuple1(t *testing.T) {
	got1 := deriveTuple1(1)()
	want1 := 1
	if got1 != want1 {
		t.Fatalf("got %d != want %d", got1, want1)
	}
}

func TestTuple2(t *testing.T) {
	got1, got2 := deriveTuple2(1, "a")()
	want1, want2 := 1, "a"
	if got1 != want1 {
		t.Fatalf("got %d != want %d", got1, want1)
	}
	if got2 != want2 {
		t.Fatalf("got %s != want %s", got2, want2)
	}
}

func TestTuple3(t *testing.T) {
	b := random(&BuiltInTypes{}).(*BuiltInTypes)
	got1, got2, got3 := deriveTuple3(1, "a", b)()
	want1, want2, want3 := 1, "a", b
	if got1 != want1 {
		t.Fatalf("got %d != want %d", got1, want1)
	}
	if got2 != want2 {
		t.Fatalf("got %s != want %s", got2, want2)
	}
	if !want3.Equal(got3) {
		t.Fatalf("got %v != want %v", got3, want3)
	}
}

type reader struct{}

func (*reader) Read(p []byte) (n int, err error) {
	return 0, errors.New("abc")
}

func TestTupleError(t *testing.T) {
	d, err := ioutil.ReadAll(&reader{})
	tup := deriveTupleError(d, err)
	_, got := tup()
	if got == nil {
		t.Fatal("expected error")
	}
}

func TestTupleReturnError(t *testing.T) {
	tup := deriveTupleError(ioutil.ReadAll(&reader{}))
	_, got := tup()
	if got == nil {
		t.Fatal("expected error")
	}
}
