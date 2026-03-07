// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

// URL type constants for DetectURLType.
const (
	// UnsupportedURLType is returned for empty or unrecognized URLs.
	UnsupportedURLType = ""

	// FileURL indicates a file URL (file://).
	FileURL = "file"

	// HttpURL indicates an HTTP or HTTPS URL.
	HttpURL = "http"

	// GitURL indicates a URL ending in .git (e.g. repository clone URL).
	GitURL = "git"
)

// DetectURLType returns the URL type (FileURL, HttpURL, GitURL, or UnsupportedURLType).
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
