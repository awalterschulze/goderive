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

func TestDup(t *testing.T) {
	c := make(chan int)
	c1, c2 := deriveDup(c)
	out := make(chan int, 100)
	done1, done2 := make(chan struct{}), make(chan struct{})
	go func() {
		for v := range c1 {
			out <- v
		}
		close(done1)
	}()
	go func() {
		for v := range c2 {
			out <- v
		}
		close(done2)
	}()
	c <- 1
	c <- 2
	close(c)
	<-done1
	<-done2
	close(out)
	want := 0
	for v := range out {
		want += v
	}
	got := 6
	if got != want {
		t.Fatalf("got %d != want %d", got, want)
	}
}
