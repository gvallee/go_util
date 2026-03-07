//
// Copyright (c) 2021-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package timestamp

import (
	"regexp"
	"testing"
)

// Now() returns a string in format YYMMDDHHMMSS (e.g. 250307143052 for 2025-03-07 14:30:52)
var nowFormat = regexp.MustCompile(`^\d{12}$`)

func TestNow(t *testing.T) {
	s := Now()
	if len(s) != 12 {
		t.Errorf("Now() length = %d, want 12", len(s))
	}
	if !nowFormat.MatchString(s) {
		t.Errorf("Now() = %q, want 12 digits (YYMMDDHHMMSS)", s)
	}
}
