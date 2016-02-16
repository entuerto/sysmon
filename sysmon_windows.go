// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysmon

import (
	"time"

	"github.com/entuerto/sysmon/internal/win32"
)

func upTime() time.Duration {
	d := win32.GetTickCount64()
	return time.Duration(d) * time.Millisecond
}