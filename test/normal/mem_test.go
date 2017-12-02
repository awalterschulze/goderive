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

import "testing"
import "fmt"

type Adder struct {
	Int int
}

func TestMemGet(t *testing.T) {
	called := 0
	get := func() *BuiltInTypes {
		called++
		return &BuiltInTypes{Int: 1}
	}
	mget := deriveMemGet(get)
	if got, want := mget(), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := mget(), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
}

func TestMemInc(t *testing.T) {
	called := 0
	inc := func(n int) int {
		called++
		return n + 1
	}
	minc := deriveMemInc(inc)
	if got, want := minc(1), 2; want != got {
		t.Fatalf("inc(1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := minc(1), 2; want != got {
		t.Fatalf("inc(1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := minc(2), 3; want != got {
		t.Fatalf("inc(2) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := minc(1), 2; want != got {
		t.Fatalf("inc(1) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := minc(2), 3; want != got {
		t.Fatalf("inc(2) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
}

func TestMemIncTo(t *testing.T) {
	called := 0
	inc := func(a Adder) int {
		called++
		return a.Int + 1
	}
	minc := deriveMemIncTo(inc)
	if got, want := minc(Adder{Int: 1}), 2; want != got {
		t.Fatalf("inc(Adder{Int: 1}) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := minc(Adder{Int: 1}), 2; want != got {
		t.Fatalf("inc(Adder{Int: 1}) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := minc(Adder{Int: 2}), 3; want != got {
		t.Fatalf("inc(Adder{Int: 2}) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := minc(Adder{Int: 1}), 2; want != got {
		t.Fatalf("inc(Adder{Int: 1}) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := minc(Adder{Int: 2}), 3; want != got {
		t.Fatalf("inc(Adder{Int: 2}) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
}

func TestMemAdd(t *testing.T) {
	called := 0
	add := func(a, b int) int {
		called++
		return a + b
	}
	madd := deriveMemAdd(add)
	if got, want := madd(1, 1), 2; want != got {
		t.Fatalf("add(1, 1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := madd(1, 1), 2; want != got {
		t.Fatalf("add(1, 1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := madd(1, 2), 3; want != got {
		t.Fatalf("add(1, 2) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := madd(1, 1), 2; want != got {
		t.Fatalf("add(1, 1) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := madd(2, 1), 3; want != got {
		t.Fatalf("add(2, 1) got %d want %d", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
	if got, want := madd(1, 1), 2; want != got {
		t.Fatalf("add(1, 1) got %d want %d", got, want)
	}
	if got, want := madd(1, 2), 3; want != got {
		t.Fatalf("add(1, 2) got %d want %d", got, want)
	}
	if got, want := madd(2, 1), 3; want != got {
		t.Fatalf("add(2, 1) got %d want %d", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
}

func TestMemAddTo(t *testing.T) {
	called := 0
	add := func(a Adder, b int) int {
		called++
		return a.Int + b
	}
	madd := deriveMemAddTo(add)
	if got, want := madd(Adder{Int: 1}, 1), 2; want != got {
		t.Fatalf("add(Adder{Int: 1}, 1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := madd(Adder{Int: 1}, 1), 2; want != got {
		t.Fatalf("add(Adder{Int: 1}, 1) got %d want %d", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := madd(Adder{Int: 1}, 2), 3; want != got {
		t.Fatalf("add(Adder{Int: 1}, 2) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := madd(Adder{Int: 1}, 1), 2; want != got {
		t.Fatalf("add(Adder{Int: 1}, 1) got %d want %d", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := madd(Adder{Int: 2}, 1), 3; want != got {
		t.Fatalf("add(Adder{Int: 2}, 1) got %d want %d", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
	if got, want := madd(Adder{Int: 1}, 1), 2; want != got {
		t.Fatalf("add(Adder{Int: 1}, 1) got %d want %d", got, want)
	}
	if got, want := madd(Adder{Int: 1}, 2), 3; want != got {
		t.Fatalf("add(Adder{Int: 1}, 2) got %d want %d", got, want)
	}
	if got, want := madd(Adder{Int: 2}, 1), 3; want != got {
		t.Fatalf("add(Adder{Int: 2}, 1) got %d want %d", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
}

func TestMemSet(t *testing.T) {
	called := 0
	set := func(a *BuiltInTypes, b int) *BuiltInTypes {
		called++
		c := a.Clone()
		c.Int = b
		return c
	}
	mset := deriveMemSet(set)
	if got, want := mset(&BuiltInTypes{Int: 1}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := mset(&BuiltInTypes{Int: 1}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	if got, want := mset(&BuiltInTypes{Int: 1}, 2), (&BuiltInTypes{Int: 2}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 2) got %v want %v", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := mset(&BuiltInTypes{Int: 1}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	if got, want := mset(&BuiltInTypes{Int: 2}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 2}, 1) got %v want %v", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
	if got, want := mset(&BuiltInTypes{Int: 1}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if got, want := mset(&BuiltInTypes{Int: 1}, 2), (&BuiltInTypes{Int: 2}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 1}, 2) got %v want %v", got, want)
	}
	if got, want := mset(&BuiltInTypes{Int: 2}, 1), (&BuiltInTypes{Int: 1}); !want.Equal(got) {
		t.Fatalf("set(&BuiltInTypes{Int: 2}, 1) got %v want %v", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
}

func TestMemSetError(t *testing.T) {
	called := 0
	seterr := func(a *BuiltInTypes, b int) (*BuiltInTypes, error) {
		called++
		c := a.Clone()
		c.Int = b
		if b == 2 {
			return nil, fmt.Errorf("error b == 2")
		}
		return c, nil
	}
	mseterr := deriveMemSetErr(seterr)
	got, goterr := mseterr(&BuiltInTypes{Int: 1}, 1)
	want, wanterr := (&BuiltInTypes{Int: 1}), error(nil)
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 1}, 1)
	want, wanterr = &BuiltInTypes{Int: 1}, nil
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 1 {
		t.Fatalf("not called once, but %d", called)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 1}, 2)
	want, wanterr = nil, fmt.Errorf("error b = 2")
	if goterr == nil || wanterr == nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 2) got %v want %v", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 1}, 1)
	want, wanterr = &BuiltInTypes{Int: 1}, nil
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	if called != 2 {
		t.Fatalf("not called twice, but %d", called)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 2}, 1)
	want, wanterr = &BuiltInTypes{Int: 1}, nil
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 2}, 1) got %v want %v", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 1}, 1)
	want, wanterr = &BuiltInTypes{Int: 1}, nil
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 1) got %v want %v", got, want)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 1}, 2)
	want, wanterr = nil, fmt.Errorf("error b = 2")
	if goterr == nil || wanterr == nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 1}, 2) got %v want %v", got, want)
	}
	got, goterr = mseterr(&BuiltInTypes{Int: 2}, 1)
	want, wanterr = &BuiltInTypes{Int: 1}, nil
	if goterr != nil || wanterr != nil || !want.Equal(got) {
		t.Fatalf("seterr(&BuiltInTypes{Int: 2}, 1) got %v want %v", got, want)
	}
	if called != 3 {
		t.Fatalf("not called thrice, but %d", called)
	}
}
