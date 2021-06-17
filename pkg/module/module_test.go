//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package module

import "testing"

func TestToEnv(t *testing.T) {
	ToEnv([]string{"PATH"}, []string{"intel/2019.5"})
}
