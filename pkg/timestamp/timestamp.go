//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package timestamp

import "time"

func Now() string {
	now := time.Now()
	return string(now.Format("060102150405"))
}
