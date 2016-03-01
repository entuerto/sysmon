// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"fmt"
	"log"
	"testing"
	"time"
)


func ExampleSwapMemory() {
	s, err := SwapMemory()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", s)
	// Output: _
}

func ExampleVirtualMemory() {
	v, err := VirtualMemory()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", v)
	// Output: _
}

func TestMemoryPerformance(t *testing.T) {
	qpi, err := QueryPerformanceInformation(time.Second)

	if err != nil {
		t.Errorf("error %v", err)
	}

	fmt.Println(" System   Commit                        Physical             Kernel                         Count")
	fmt.Println(" Cache    Total     Limit     Peak      Total     Available  Total     Paged     Nonpaged   Handle  Process  Thread")
	go func() {
		for {
			pc := <- qpi.PerfCounterChan
			fmt.Printf(" %7s %9s %9s %9s %9s %9s %11s %9s %9s %7d %8d %7d\n", 
				pc.SystemCache, 
		    	pc.CommitTotal,
				pc.CommitLimit,      
				pc.CommitPeak,       
				pc.PhysicalTotal,
				pc.PhysicalAvailable,
				pc.KernelTotal,
				pc.KernelPaged,
				pc.KernelNonpaged,
				pc.HandleCount,
				pc.ProcessCount,
				pc.ThreadCount)	
		}
	}()

	<-time.After(20 * time.Second)
    qpi.Stop()
}
      