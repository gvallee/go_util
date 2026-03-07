// Copyright (c) 2019-2026, Sylabs Inc. All rights reserved.
// Copyright (c) 2026, NVIDIA CORPORATION. All rights reserved.
// This software is licensed under a 3-clause BSD license. Please consult the
// LICENSE.md file distributed with the sources of this project regarding your
// rights to use or distribute this software.

package util

import (
	"testing"
)

func TestCleanupString(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"no escape", "hello world", "hello world"},
		{"literal backslash x1b reset", "foo\\x1b[0mbar", "foobar"},
		{"literal backslash x1b color", "foo\\x1b[33mbar", "foobar"},
		{"actual ANSI reset", "foo\x1b[0mbar", "foobar"},
		{"actual ANSI yellow", "foo\x1b[33mbar", "foobar"},
		{"actual ANSI multiple", "\x1b[1m\x1b[33mbold yellow\x1b[0m", "bold yellow"},
		{"mixed digits", "\x1b[31mred\x1b[0m", "red"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanupString(tt.in)
			if got != tt.want {
				t.Errorf("CleanupString(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
