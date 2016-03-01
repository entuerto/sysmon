// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sysmon

import (
	"fmt"
	"time"
)

type Size uint64

const (
	_       = iota // ignore first value by assigning to blank identifier
	KB Size = 1 << (10 * iota)
	MB
	GB
	TB
	PB
	EB
//	ZB
//	YB
)

func (b Size) String() string {
	switch {
/*		
	case b >= YB:
	    return fmt.Sprintf("%.2f YB", float64(b) / float64(YB))
	case b >= ZB:
	    return fmt.Sprintf("%.2f ZB", float64(b) / float64(ZB))
*/	    
	case b >= EB:
	    return fmt.Sprintf("%.2f EB", float64(b) / float64(EB))
	case b >= PB:
	    return fmt.Sprintf("%.2f PB", float64(b) / float64(PB))
	case b >= TB:
	    return fmt.Sprintf("%.2f TB", float64(b) / float64(TB))
	case b >= GB:
	    return fmt.Sprintf("%.2f GB", float64(b) / float64(GB))
	case b >= MB:
	    return fmt.Sprintf("%.2f MB", float64(b) / float64(MB))
	case b >= KB:
	    return fmt.Sprintf("%.2f KB", float64(b) / float64(KB))
	}
	return fmt.Sprintf("%.2d B", b)
}

func UpTime() time.Duration {
	return upTime()
}