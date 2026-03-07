//
// Copyright (c) 2020-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package notation

import (
	"testing"
)

func TestCompressIntArray(t *testing.T) {
	tests := []struct {
		array           []int
		expectedResults string
	}{
		{[]int{0, 1, 2, 3, 4, 5, 6, 8, 9, 10, 42}, "0-6,8-10,42"},
		{[]int{}, ""},
		{[]int{7}, "7"},
		{[]int{1, 2, 3}, "1-3"},
	}

	for _, tt := range tests {
		str := CompressIntArray(tt.array)
		if str != tt.expectedResults {
			t.Errorf("CompressIntArray(%v) = %q, want %q", tt.array, str, tt.expectedResults)
		}
	}
}

func TestConvertCompressedCallListToIntSlice(t *testing.T) {
	tests := []struct {
		input         string
		expectedArray []int
	}{
		{"32", []int{32}},
		{"32-35,40", []int{32, 33, 34, 35, 40}},
		{"", []int{}},
		{"1, 2, 3", []int{1, 2, 3}},
	}

	for _, tt := range tests {
		array, err := ConvertCompressedCallListToIntSlice(tt.input)
		if err != nil {
			t.Fatalf("ConvertCompressedCallListToIntSlice(%q) failed: %v", tt.input, err)
		}
		if len(array) != len(tt.expectedArray) {
			t.Fatalf("ConvertCompressedCallListToIntSlice(%q) returned %d elements, want %d", tt.input, len(array), len(tt.expectedArray))
		}
		for i := range tt.expectedArray {
			if array[i] != tt.expectedArray[i] {
				t.Fatalf("element %d = %d, want %d", i, array[i], tt.expectedArray[i])
			}
		}
	}
}

func TestConvertCompressedCallListToIntSliceErrors(t *testing.T) {
	_, err := ConvertCompressedCallListToIntSlice("abc")
	if err == nil {
		t.Error("ConvertCompressedCallListToIntSlice(\"abc\") expected error")
	}
	_, err = ConvertCompressedCallListToIntSlice("1-2a")
	if err == nil {
		t.Error("ConvertCompressedCallListToIntSlice(\"1-2a\") expected error")
	}
}

func TestGetNumberOfRanksFromCompressedNotation(t *testing.T) {
	tests := []struct {
		input          string
		expectedOutput int
	}{
		{"1, 2", 2},
		{"1,2", 2},
		{"1-5", 5},
		{"0,1-5", 6},
		{"0,1-5,6", 7},
		{"32", 1},
		{"", 0},
	}
	for _, tt := range tests {
		n, err := GetNumberOfRanksFromCompressedNotation(tt.input)
		if err != nil {
			t.Fatalf("GetNumberOfRanksFromCompressedNotation(%q) failed: %v", tt.input, err)
		}
		if n != tt.expectedOutput {
			t.Fatalf("GetNumberOfRanksFromCompressedNotation(%q) = %d, want %d", tt.input, n, tt.expectedOutput)
		}
	}
}

func TestGetNumberOfEltsFromCompressedNotationErrors(t *testing.T) {
	_, err := GetNumberOfEltsFromCompressedNotation("1-x")
	if err == nil {
		t.Error("GetNumberOfEltsFromCompressedNotation(\"1-x\") expected error")
	}
	_, err = GetNumberOfEltsFromCompressedNotation("a")
	if err == nil {
		t.Error("GetNumberOfEltsFromCompressedNotation(\"a\") expected error")
	}
}

func TestIntSliceToString(t *testing.T) {
	tests := []struct {
		s              []int
		expectedResult string
	}{
		{[]int{1, 2, 3, 4, 5, 6}, "1,2,3,4,5,6"},
		{[]int{}, ""},
		{[]int{42}, "42"},
	}

	for _, tt := range tests {
		str := IntSliceToString(tt.s)
		if str != tt.expectedResult {
			t.Errorf("IntSliceToString(%v) = %q, want %q", tt.s, str, tt.expectedResult)
		}
	}
}

func TestIntToAA(t *testing.T) {
	tests := []struct {
		input          int
		expectedOutput string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
		{78, "CA"},
	}

	for _, tt := range tests {
		res := IntToAA(tt.input)
		if res != tt.expectedOutput {
			t.Errorf("IntToAA(%d) = %q, want %q", tt.input, res, tt.expectedOutput)
		}
	}
}

func TestIntToAANegative(t *testing.T) {
	if got := IntToAA(-1); got != "-ERROR-" {
		t.Errorf("IntToAA(-1) = %q, want \"-ERROR-\"", got)
	}
}
