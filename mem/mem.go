// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"fmt"
	"strings"

	"github.com/entuerto/sysmon"
)

type Swap struct {
	Total   sysmon.Size `json:"total"`   // total swap memory in bytes
	Used    sysmon.Size `json:"used"`    // used swap memory in bytes
	Free    sysmon.Size `json:"free"`    // free swap memory in bytes
	Percent float64     `json:"percent"` // the percentage usage calculated as (total - available) / total * 100
	SIn     sysmon.Size `json:"sin"`     // the number of bytes the system has swapped in from disk (cumulative)
	SOut    sysmon.Size `json:"sout"`    // the number of bytes the system has swapped out from disk (cumulative)
}

func (sw Swap) GoString() string {
	s := []string{"Swap{", 
			fmt.Sprintf("  Total   : %s", sw.Total), 
			fmt.Sprintf("  Used    : %s", sw.Used), 
			fmt.Sprintf("  Free    : %s", sw.Free), 
			fmt.Sprintf("  Percent : %.2f", sw.Percent), 
			fmt.Sprintf("  SIn     : %s", sw.SIn), 
			fmt.Sprintf("  SOut    : %s", sw.SOut), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func SwapMemory() (*Swap, error) {
	return swapMemory()
}

//---------------------------------------------------------------------------------------

type Virtual struct {
	Total     sysmon.Size `json:"total"`
	Available sysmon.Size `json:"available"`
	Used      sysmon.Size `json:"used"`
	Free      sysmon.Size `json:"free"`
	Percent   float64     `json:"percent"` 
}

func (v Virtual) GoString() string {
	s := []string{"Virtual{", 
			fmt.Sprintf("  Total     : %s", v.Total), 
			fmt.Sprintf("  Available : %s", v.Available), 
			fmt.Sprintf("  Used      : %s", v.Used), 
			fmt.Sprintf("  Free      : %s", v.Free), 
			fmt.Sprintf("  Percent   : %.2f", v.Percent), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func VirtualMemory() (*Virtual, error) {
	return virtualMemory()
}