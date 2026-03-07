//
// Copyright (c) 2021-2026, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

// Package timestamp provides simple time formatting for filenames and logs.
package timestamp

import "time"

// Now returns the current time as YYMMDDHHMMSS (e.g. 260307143052 for 2026-03-07 14:30:52).
func Now() string {
	return time.Now().Format("060102150405")
}
