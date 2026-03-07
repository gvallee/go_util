//
// Copyright (c) 2020-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

// Package notation provides compressed integer list notation (e.g. "1-5,8,10") and helpers.
package notation

import (
	"fmt"
	"strconv"
	"strings"
)

func addRange(str string, start int, end int) string {
	if str == "" {
		return fmt.Sprintf("%d-%d", start, end)
	}
	return fmt.Sprintf("%s,%d-%d", str, start, end)
}

func addSingleton(str string, n int) string {

	if str == "" {
		return fmt.Sprintf("%d", n)
	}

	return fmt.Sprintf("%s,%d", str, n)
}

// CompressIntArray compresses consecutive integers into range notation (e.g. [1,2,3] -> "1-3").
// Singletons and ranges are comma-separated; input need not be sorted.
func CompressIntArray(array []int) string {
	compressedRep := ""
	for i := 0; i < len(array); i++ {
		start := i
		for i+1 < len(array) && array[i]+1 == array[i+1] {
			i++
		}
		if i != start {
			// We found a range
			compressedRep = addRange(compressedRep, array[start], array[i])
		} else {
			// We found a singleton
			compressedRep = addSingleton(compressedRep, array[i])
		}
	}
	return compressedRep
}

// GetNumberOfEltsFromCompressedNotation parses a compressed notation string (e.g. "1-5,8")
// and returns the total number of elements. Accepts ", " or "," as separator; empty string returns 0.
func GetNumberOfEltsFromCompressedNotation(str string) (int, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0, nil
	}
	num := 0
	// ", " and "," are two possible delimiters but cannot be mixed
	t1 := strings.Split(str, ", ")
	if len(t1) == 1 {
		t1 = strings.Split(str, ",")
	}
	for _, t := range t1 {
		t2 := strings.Split(t, "-")
		if len(t2) == 2 {
			val1, err := strconv.Atoi(strings.TrimSpace(t2[0]))
			if err != nil {
				return 0, err
			}
			val2, err := strconv.Atoi(strings.TrimSpace(t2[1]))
			if err != nil {
				return 0, err
			}
			if val1 > val2 {
				return 0, fmt.Errorf("invalid range %d-%d", val1, val2)
			}
			num += val2 - val1 + 1
		} else {
			for _, part := range t2 {
				_, err := strconv.Atoi(strings.TrimSpace(part))
				if err != nil {
					return 0, err
				}
				num++
			}
		}
	}
	return num, nil
}

// GetNumberOfRanksFromCompressedNotation returns the number of elements in the compressed notation;
// it is equivalent to GetNumberOfEltsFromCompressedNotation.
func GetNumberOfRanksFromCompressedNotation(str string) (int, error) {
	return GetNumberOfEltsFromCompressedNotation(str)
}

// ConvertCompressedCallListToIntSlice parses a compressed notation string (e.g. "1-5,8")
// and returns the expanded slice of integers. Empty string returns nil, nil.
func ConvertCompressedCallListToIntSlice(str string) ([]int, error) {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil, nil
	}
	var callIDs []int

	tokens := strings.Split(str, ", ")
	tokensNoSpace := strings.Split(str, ",")
	if len(tokens) == 1 && len(tokensNoSpace) > 1 {
		//log.Printf("[INFO] Switching to notation without spaces for %s", str)
		tokens = tokensNoSpace
	}

	for _, t := range tokens {
		tokens2 := strings.Split(t, "-")
		if len(tokens2) == 2 {
			val1, err := strconv.Atoi(strings.TrimSpace(tokens2[0]))
			if err != nil {
				return callIDs, err
			}
			val2, err := strconv.Atoi(strings.TrimSpace(tokens2[1]))
			if err != nil {
				return callIDs, err
			}
			if val1 > val2 {
				return callIDs, fmt.Errorf("invalid range %d-%d", val1, val2)
			}
			for i := val1; i <= val2; i++ {
				callIDs = append(callIDs, i)
			}
		} else {
			for _, t2 := range tokens2 {
				n, err := strconv.Atoi(strings.TrimSpace(t2))
				if err != nil {
					return callIDs, fmt.Errorf("unable to parse compressed string: %s", str)
				}
				callIDs = append(callIDs, n)
			}
		}
	}

	//log.Printf("%s was decompressed into %d elements\n", str, len(callIDs))

	return callIDs, nil
}

// IntSliceToString formats the slice as a comma-separated list (e.g. "1,2,3").
func IntSliceToString(s []int) string {
	if len(s) == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(len(s) * 4) // rough: 3 digits + comma per element
	b.WriteString(strconv.Itoa(s[0]))
	for i := 1; i < len(s); i++ {
		b.WriteByte(',')
		b.WriteString(strconv.Itoa(s[i]))
	}
	return b.String()
}

var alphabet = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

// IntToAA converts a non-negative integer to alphabetic notation (0=A, 26=AA, 27=AB, ...).
// Returns "-ERROR-" for negative input.
func IntToAA(num int) string {
	if num < 0 {
		return "-ERROR-"
	}
	if num >= 0 && num < 26 {
		return alphabet[num]
	}

	idx := num % 26
	return IntToAA(num/26-1) + alphabet[idx]
}
