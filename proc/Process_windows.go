// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"
//	"log"
	"os"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/entuerto/sysmon/internal/win32"
	"github.com/entuerto/sysmon"
)

func OpenProcess(pid uint32) (*Process, error) {
	const da = syscall.STANDARD_RIGHTS_READ | syscall.PROCESS_QUERY_INFORMATION | syscall.SYNCHRONIZE

	h, err := syscall.OpenProcess(da, false, uint32(pid))
	if err != nil {
		return nil, os.NewSyscallError("OpenProcess", err)
	}

	procEntry, err := win32.FindProcessEntry(pid) 
	if err != nil {
		return nil, err
	}

	modEntry, err := win32.FirstModuleEntry(pid)
	if err != nil {
		return nil, err
	}

	handleCount, err := win32.GetProcessHandleCount(h)
	if err != nil {
		return nil, err
	}

	return &Process{
		Pid        : pid,
		handle     : uintptr(h),
		ParentId   : procEntry.ParentProcessID,
		Name       : syscall.UTF16ToString(modEntry.ModuleName[:]),
		Executable : syscall.UTF16ToString(procEntry.ExeFile[:]),
		CmdLine    : syscall.UTF16ToString(modEntry.ExePath[:]),
		HandleCount: handleCount,
		ThreadCount: procEntry.Threads,
	}, nil
}

func (p Process) ioCounters() (*IOCounters, error){
	wioc, err := win32.GetProcessIoCounters(syscall.Handle(p.handle)) 
	if err != nil {
		return nil, err
	}

	return &IOCounters{
		ReadCount  : wioc.ReadOperationCount,
		WriteCount : wioc.WriteOperationCount,
		ReadBytes  : sysmon.Size(wioc.ReadTransferCount),
		WriteBytes : sysmon.Size(wioc.WriteTransferCount),
	}, nil
}

func (p Process) usage() (*TimeUsage, error) {
	var u syscall.Rusage

	err := syscall.GetProcessTimes(syscall.Handle(p.handle), &u.CreationTime, &u.ExitTime, &u.KernelTime, &u.UserTime)
	if err != nil {
		return nil, os.NewSyscallError("GetProcessTimes", err)
	}
	 
	return &TimeUsage{
		CreationTime : time.Unix(0, u.CreationTime.Nanoseconds()),
		ExitTime     : time.Unix(0, u.ExitTime.Nanoseconds()),
		KernelTime   : toDuration(u.KernelTime),
		UserTime     : toDuration(u.UserTime),
	}, nil
}

type MemoryCounters struct {
	PageFaultCount             uint32      // The number of page faults.
	PeakWorkingSetSize         sysmon.Size // The peak working set size, in bytes.
	WorkingSetSize             sysmon.Size // The current working set size, in bytes.
	QuotaPeakPagedPoolUsage    sysmon.Size // The peak paged pool usage, in bytes.
	QuotaPagedPoolUsage        sysmon.Size // The current paged pool usage, in bytes.
	QuotaPeakNonPagedPoolUsage sysmon.Size // The peak nonpaged pool usage, in bytes.
	QuotaNonPagedPoolUsage     sysmon.Size // The current nonpaged pool usage, in bytes.
	PagefileUsage              sysmon.Size // The Commit Charge value in bytes for this process. Commit Charge 
	                                       // is the total amount of memory that the memory manager has committed 
	                                       // for a running process.
	PeakPagefileUsage          sysmon.Size // The peak value in bytes of the Commit Charge during the lifetime 
	                                       // of this process.
}

func (mc MemoryCounters) GoString() string {
	s := []string{"MemoryCounters{", 
			fmt.Sprintf("  PageFaultCount             : %d", mc.PageFaultCount),   
			fmt.Sprintf("  PeakWorkingSetSize         : %s", mc.PeakWorkingSetSize),   
			fmt.Sprintf("  WorkingSetSize             : %s", mc.WorkingSetSize),   
			fmt.Sprintf("  QuotaPeakPagedPoolUsage    : %s", mc.QuotaPeakPagedPoolUsage),   
			fmt.Sprintf("  QuotaPagedPoolUsage        : %s", mc.QuotaPagedPoolUsage),   
			fmt.Sprintf("  QuotaPeakNonPagedPoolUsage : %s", mc.QuotaPeakNonPagedPoolUsage),   
			fmt.Sprintf("  QuotaNonPagedPoolUsage     : %s", mc.QuotaNonPagedPoolUsage),   
			fmt.Sprintf("  PagefileUsage              : %s", mc.PagefileUsage),   
			fmt.Sprintf("  PeakPagefileUsage          : %s", mc.PeakPagefileUsage),   
			"}",
	}
	return strings.Join(s, "\n")	
}

func (p Process) memoryInfo() (*MemoryCounters, error) {
	pmc, err := win32.GetProcessMemoryInfo(syscall.Handle(p.handle)) 
	if err != nil {
		return nil, err
	} 

	return &MemoryCounters{
		PageFaultCount             : pmc.PageFaultCount,
		PeakWorkingSetSize         : sysmon.Size(pmc.PeakWorkingSetSize),
		WorkingSetSize             : sysmon.Size(pmc.WorkingSetSize),
		QuotaPeakPagedPoolUsage    : sysmon.Size(pmc.QuotaPeakPagedPoolUsage),
		QuotaPagedPoolUsage        : sysmon.Size(pmc.QuotaPagedPoolUsage),
		QuotaPeakNonPagedPoolUsage : sysmon.Size(pmc.QuotaPeakNonPagedPoolUsage),
		QuotaNonPagedPoolUsage     : sysmon.Size(pmc.QuotaNonPagedPoolUsage),
		PagefileUsage              : sysmon.Size(pmc.PagefileUsage),
		PeakPagefileUsage          : sysmon.Size(pmc.PeakPagefileUsage),
	}, nil
}

type Module struct {
	ProcessID uint32
	BaseAddr  uintptr      // The base address of the module in the context of the owning process.
	BaseSize  sysmon.Size  // The size of the module, in bytes.
	Handle    uintptr      // A handle to the module in the context of the owning process.
	Name      string
	ExePath   string
}

func (p Process) modules() ([]*Module, error) {
	snapshot, err := win32.CreateToolhelp32Snapshot(win32.TH32CS_SNAPMODULE, p.Pid)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	var modEntry win32.ModuleEntry32
	modEntry.Size = uint32(unsafe.Sizeof(modEntry))

	if err := win32.Module32First(snapshot, &modEntry); err != nil {
		return nil, err 
	}

	for {
	//	fmt.Println("ModuleID:     ", modEntry.ModuleID)
	//	fmt.Println("ModuleName:   ", syscall.UTF16ToString(modEntry.ModuleName[:]))
		fmt.Println("ExePath:      ", syscall.UTF16ToString(modEntry.ExePath[:]))

		err = win32.Module32Next(snapshot, &modEntry)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func toDuration(ft syscall.Filetime) time.Duration {
	n := int64(ft.HighDateTime) << 32 + int64(ft.LowDateTime) // in 100-nanosecond intervals
	return time.Duration(n * 100) * time.Nanosecond
}