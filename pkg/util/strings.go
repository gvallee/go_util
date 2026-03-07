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

// CleanupString cleans up a string by removing ANSI color escape sequences.
// This is convenient when a benchmark returns colored text to the terminal
// that we then try to parse. Handles both actual ESC bytes (0x1b) and
// literal \x1b character sequences.
func CleanupString(str string) string {
	// Remove ANSI escape sequences: ESC [ <digits> m (actual 0x1b byte)
	ansiEscape := regexp.MustCompile("\x1b\\[[0-9]+m")
	str = ansiEscape.ReplaceAllString(str, "")

	// Remove literal \x1b sequences (e.g. from escaped log output)
	literalEscape := regexp.MustCompile(`\\x1b\[[0-9]+m`)
	str = literalEscape.ReplaceAllString(str, "")

	str = strings.ReplaceAll(str, "\x1b[0m", "")
	str = strings.ReplaceAll(str, "\x1b[33m", "")
	str = strings.ReplaceAll(str, `\x1b`+"[0m", "")
	return strings.ReplaceAll(str, `\x1b`+"[33m", "")
}
