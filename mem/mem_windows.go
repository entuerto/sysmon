// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/entuerto/sysmon"
	"github.com/entuerto/sysmon/internal/win32"
)

func swapMemory() (*Swap, error) {
	mem, err := win32.GlobalMemoryStatusEx() 
	if err != nil {
		return nil, err
	}

	return &Swap{
		Total   : sysmon.Size(mem.TotalPageFile),
		Used    : sysmon.Size(mem.TotalPageFile - mem.AvailPageFile),
		Free    : sysmon.Size(mem.AvailPageFile), 
		Percent : float64(mem.TotalPageFile - mem.AvailPageFile) / float64(mem.TotalPageFile) * 100, // (total - available) / total * 100
	}, nil
}

func virtualMemory() (*Virtual, error) {
	mem, err := win32.GlobalMemoryStatusEx() 
	if err != nil {
		return nil, err
	}

	return &Virtual{
		Total     : sysmon.Size(mem.TotalPhys),
		Available : sysmon.Size(mem.AvailPhys),
		Used      : sysmon.Size(mem.AvailPhys), 
		Free      : sysmon.Size(mem.TotalPhys - mem.AvailPhys), 
		Percent   : float64(mem.TotalPhys - mem.AvailPhys) / float64(mem.TotalPhys) * 100,
	}, nil
}

type PerfCounter struct {
	CommitTotal       sysmon.Size // Committed Bytes - System: Page File first value (in MB)
	CommitLimit       sysmon.Size // Commit Limit    - System: Page File second value (in MB)
	CommitPeak        sysmon.Size // Commit Charge: Peak
	                  
	PhysicalTotal     sysmon.Size // Physical Memory - Total
	PhysicalAvailable sysmon.Size // Physical Memory - Available KB

	SystemCache       sysmon.Size // Cache Bytes + Sharable pages on the standby and modified lists
	                  
	KernelTotal       sysmon.Size // Kernel Memory: Total - Pool Paged Bytes + Pool Nonpaged Bytes
	KernelPaged       sysmon.Size // Kernel Memory: Paged - Pool Paged Bytes
	KernelNonpaged    sysmon.Size // Kernel Memory: Nonpaged - Pool Nonpaged Bytes

	HandleCount       uint32 // The current number of open handles.
	ProcessCount      uint32 // The current number of processes.
	ThreadCount       uint32 // The current number of threads.     

}

func (pc PerfCounter) GoString() string {
	s := []string{"PerfCounter{", 
			fmt.Sprintf("  CommitTotal       : %s", pc.CommitTotal), 
			fmt.Sprintf("  CommitLimit       : %s", pc.CommitLimit), 
			fmt.Sprintf("  CommitPeak        : %s", pc.CommitPeak), 
			fmt.Sprintf("  PhysicalTotal     : %s", pc.PhysicalTotal), 
			fmt.Sprintf("  PhysicalAvailable : %s", pc.PhysicalAvailable), 
			fmt.Sprintf("  SystemCache       : %s", pc.SystemCache), 
			fmt.Sprintf("  KernelTotal       : %s", pc.KernelTotal), 
			fmt.Sprintf("  KernelPaged       : %s", pc.KernelPaged), 
			fmt.Sprintf("  KernelNonpaged    : %s", pc.KernelNonpaged), 
			fmt.Sprintf("  HandleCount       : %d", pc.HandleCount), 
			fmt.Sprintf("  ProcessCount      : %d", pc.ProcessCount), 
			fmt.Sprintf("  ThreadCount       : %d", pc.ThreadCount), 
			"}",
	}
	return strings.Join(s, "\n")	
}

type QueryPerfInfo struct {
	PerfCounterChan chan *PerfCounter
	quit            chan bool
}

// Stop signals the goroutine to stop querying PerfCounters.
func (q QueryPerfInfo) Stop() {
	go func() {
		q.quit <- true
	}()
}

// System Memory Performance Information
func QueryPerformanceInformation(freq time.Duration) (QueryPerfInfo, error) {
	qpi := QueryPerfInfo{
		PerfCounterChan : make(chan *PerfCounter),
		quit            : make(chan bool),
	}

	go queryPerformanceInformation(freq, qpi)

	return qpi, nil 
}

// 
func queryPerformanceInformation(freq time.Duration, qpi QueryPerfInfo) {

	for {
		select {
		case <- time.After(freq):
			pi, err := win32.GetPerformanceInfo() 
		
			if err != nil {
				log.Fatal(err)
			}

			pc := &PerfCounter{
				CommitTotal       : sysmon.Size(pi.CommitTotal * pi.PageSize),
				CommitLimit       : sysmon.Size(pi.CommitLimit * pi.PageSize),
				CommitPeak        : sysmon.Size(pi.CommitPeak * pi.PageSize),
				PhysicalTotal     : sysmon.Size(pi.PhysicalTotal * pi.PageSize),
				PhysicalAvailable : sysmon.Size(pi.PhysicalAvailable * pi.PageSize),
				SystemCache       : sysmon.Size(pi.SystemCache * pi.PageSize),
				KernelTotal       : sysmon.Size(pi.KernelTotal * pi.PageSize),
				KernelPaged       : sysmon.Size(pi.KernelPaged * pi.PageSize),
				KernelNonpaged    : sysmon.Size(pi.KernelNonpaged * pi.PageSize),
				HandleCount       : pi.HandleCount,
				ProcessCount      : pi.ProcessCount,
				ThreadCount       : pi.ThreadCount,
			}
			qpi.PerfCounterChan <- pc
		case <- qpi.quit:
			return
		}
	}
}