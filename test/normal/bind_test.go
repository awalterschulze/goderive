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
	"strconv"
	"testing"
)

func TestBind(t *testing.T) {
	read := func() (string, error) {
		return "1", nil
	}
	parseFloat := func(i string) (float64, error) {
		return strconv.ParseFloat(i, 64)
	}
	got, err := deriveBind(read, parseFloat)()
	if err != nil {
		t.Fatal(err)
	}
	want := float64(1)
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}
