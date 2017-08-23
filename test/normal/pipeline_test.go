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

func TestPipeline(t *testing.T) {
	cc := derivePipeline(toChan, wordsize)
	sizes := cc(lines)
	want := 2 + 4 + 2 + 5 +
		7 + 4 + 7 + 4 +
		7 + 5 + 7 + 4 +
		7 + 7 + 7 + 4
	got := 0
	for i := range sizes {
		got += i
	}
	if got != want {
		t.Fatalf("got %d, want %d", got, want)
	}
}
