// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysmon

import (
	"fmt"
	"testing"
)

func TestUpTime(t *testing.T) {
	d := UpTime()

	if d == 0 {
		t.Error("error zero value for UpTime()")
	}

	fmt.Printf("Uptime: %+v \n\n", d)
}
