// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

// Constants defining the URL types
const (
	// UnsupportedURLType is a constant for a type of URL that we do not support yet
	UnsupportedURLType = ""

	// FileURL is a constant for a file-based URL
	FileURL = "file"

	// HttpURL is a constant for a HTTP-based URL
	HttpURL = "http"

	// GitURL is a constant for a Git URL
	GitURL = "git"
)

// DetectURLType detects the type of the URL that is passed in.
// Returns UnsupportedURLType for empty or unrecognized URLs.
func DetectURLType(url string) string {
	if len(url) < 4 {
		return UnsupportedURLType
	}
	if len(url) >= 7 && url[:7] == "file://" {
		return FileURL
	}
	if len(url) >= 4 && url[len(url)-4:] == ".git" {
		return GitURL
	}
	if url[:4] == "http" {
		return HttpURL
	}
	return UnsupportedURLType
}
