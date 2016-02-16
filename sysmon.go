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
	    return fmt.Sprintf("%.2fYB", float64(b) / float64(YB))
	case b >= ZB:
	    return fmt.Sprintf("%.2fZB", float64(b) / float64(ZB))
*/	    
	case b >= EB:
	    return fmt.Sprintf("%.2fEB", float64(b) / float64(EB))
	case b >= PB:
	    return fmt.Sprintf("%.2fPB", float64(b) / float64(PB))
	case b >= TB:
	    return fmt.Sprintf("%.2fTB", float64(b) / float64(TB))
	case b >= GB:
	    return fmt.Sprintf("%.2fGB", float64(b) / float64(GB))
	case b >= MB:
	    return fmt.Sprintf("%.2fMB", float64(b) / float64(MB))
	case b >= KB:
	    return fmt.Sprintf("%.2fKB", float64(b) / float64(KB))
	}
	return fmt.Sprintf("%.2dB", b)
}

func UpTime() time.Duration {
	return upTime()
}