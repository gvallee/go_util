// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

import (
	"regexp"
	"strings"
)

var (
	ansiEscape    = regexp.MustCompile("\x1b\\[[0-9]+m")
	literalEscape = regexp.MustCompile(`\\x1b\[[0-9]+m`)
)

// CleanupString removes ANSI color escape sequences from str.
// Handles both actual ESC bytes (0x1b) and literal \x1b sequences (e.g. from escaped log output).
func CleanupString(str string) string {
	str = ansiEscape.ReplaceAllString(str, "")
	str = literalEscape.ReplaceAllString(str, "")

	str = strings.ReplaceAll(str, "\x1b[0m", "")
	str = strings.ReplaceAll(str, "\x1b[33m", "")
	str = strings.ReplaceAll(str, `\x1b`+"[0m", "")
	return strings.ReplaceAll(str, `\x1b`+"[33m", "")
}
