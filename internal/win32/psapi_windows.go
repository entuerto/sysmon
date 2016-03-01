// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
	"unsafe"
)

var (
	procGetProcessMemoryInfo = modkernel32.NewProc("K32GetProcessMemoryInfo")
	procGetPerformanceInfo   = modkernel32.NewProc("K32GetPerformanceInfo")
)

type ProcessMemoryCountersEx struct {
	cb                         uint32 // The size of the structure, in bytes.
	PageFaultCount             uint32 // The number of page faults.
	PeakWorkingSetSize         uint64 // The peak working set size, in bytes.
	WorkingSetSize             uint64 // The current working set size, in bytes.
	QuotaPeakPagedPoolUsage    uint64 // The peak paged pool usage, in bytes.
	QuotaPagedPoolUsage        uint64 // The current paged pool usage, in bytes.
	QuotaPeakNonPagedPoolUsage uint64 // The peak nonpaged pool usage, in bytes.
	QuotaNonPagedPoolUsage     uint64 // The current nonpaged pool usage, in bytes.
	PagefileUsage              uint64 // The Commit Charge value in bytes for this process. Commit Charge 
	                                  // is the total amount of memory that the memory manager has committed 
	                                  // for a running process.
	PeakPagefileUsage          uint64 // The peak value in bytes of the Commit Charge during the lifetime 
	                                  // of this process.
	PrivateUsage               uint64
}

func GetProcessMemoryInfo(h syscall.Handle) (*ProcessMemoryCountersEx, error) {
	var pmc ProcessMemoryCountersEx

	r, _, _ := procGetProcessMemoryInfo.Call(uintptr(h), 
	                                         uintptr(unsafe.Pointer(&pmc)),
	                                         uintptr(unsafe.Sizeof(pmc)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &pmc, nil
}

type PerformanceInformation struct {
	Size              uint32 // The size of this structure, in bytes.
	CommitTotal       uint64 // The number of pages currently committed by the system.
	CommitLimit       uint64 // The current maximum number of pages that can be committed
	                         // by the system without extending the paging file(s).
	CommitPeak        uint64 // The maximum number of pages that were simultaneously in 
	                         // the committed state since the last system reboot.
	PhysicalTotal     uint64 // The amount of actual physical memory, in pages.
	PhysicalAvailable uint64 // The amount of physical memory currently available, in pages. 
	SystemCache       uint64 // The amount of system cache memory, in pages. This is the 
	                         // size of the standby list plus the system working set.
	KernelTotal       uint64 // The sum of the memory currently in the paged and nonpaged 
	                         // kernel pools, in pages.
	KernelPaged       uint64 // The memory currently in the paged kernel pool, in pages.
	KernelNonpaged    uint64 // The memory currently in the nonpaged kernel pool, in pages.
	PageSize          uint64 // The size of a page, in bytes.
	HandleCount       uint32 // The current number of open handles.
	ProcessCount      uint32 // The current number of processes.
	ThreadCount       uint32 // The current number of threads.
} 

func GetPerformanceInfo() (*PerformanceInformation, error) {
	var pi PerformanceInformation

	r, _, _ := procGetPerformanceInfo.Call( 
				uintptr(unsafe.Pointer(&pi)),
				uintptr(unsafe.Sizeof(pi)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &pi, nil
}