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

func TestDoSuccess(t *testing.T) {
	f := func() (string, error) {
		return "a", nil
	}
	g := func() (int, error) {
		return 1, nil
	}
	s, i, err := deriveDo(f, g)
	if err != nil {
		t.Fatal(err)
	}
	if s != "a" || i != 1 {
		t.Fatalf("unexpected results %s %d", s, i)
	}
}

func TestDoFailure(t *testing.T) {
	f := func() (string, error) {
		return "a", fmt.Errorf("a")
	}
	g := func() (int, error) {
		return 1, nil
	}
	_, _, err := deriveDo(f, g)
	if err == nil {
		t.Fatal("expected error")
	}
}
