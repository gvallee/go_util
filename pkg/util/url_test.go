// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

import (
	"testing"
)

func TestDetectURLType(t *testing.T) {
	tests := []struct {
		url            string
		expectedResult string
	}{
		{url: "", expectedResult: UnsupportedURLType},
		{url: "a", expectedResult: UnsupportedURLType},
		{url: "ab", expectedResult: UnsupportedURLType},
		{url: "htt", expectedResult: UnsupportedURLType},
		{url: "file://aurl", expectedResult: FileURL},
		{url: "file:///path/to/file", expectedResult: FileURL},
		{url: "http://myurl", expectedResult: HttpURL},
		{url: "https://aurl", expectedResult: HttpURL},
		{url: "http://something/something.git", expectedResult: GitURL},
		{url: "https://repo.example.com/proj.git", expectedResult: GitURL},
		{url: "ftp://example.com", expectedResult: UnsupportedURLType},
		{url: "git://host/repo", expectedResult: UnsupportedURLType},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			res := DetectURLType(tt.url)
			if res != tt.expectedResult {
				t.Errorf("DetectURLType(%q) = %q, want %q", tt.url, res, tt.expectedResult)
			}
		})
	}
}
