// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"testing"
)

func TestInSlice(t *testing.T) {
	sl := []string{"A", "b"}
	if !InSlice("A", sl) {
		t.Error("should be true")
	}
	if InSlice("B", sl) {
		t.Error("should be false")
	}
}

func TestSliceEqual(t *testing.T) {
	sl := []int64{1, 2, 3}
	s2 := []int64{3, 2, 1}
	if !SliceEqual(sl, s2) {
		t.Error("should be Equal")
	}

	sl = []int64{1, 2, 3}
	s2 = []int64{3, 2}
	if SliceEqual(sl, s2) {
		t.Error("should not Equal")
	}
}

func TestSlicePadString(t *testing.T) {
	sl2 := []string{"1", "2"}
	padSlice := SlicePadString(sl2, 4, "")
	if len(padSlice) != 4 {
		t.Error("the length should be 4")
	}
}
