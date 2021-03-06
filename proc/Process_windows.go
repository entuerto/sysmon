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

//---------------------------------------------------------------------------------------

type IOCounters struct  {
	ReadCount  uint64      `json:"readCount"`
	WriteCount uint64      `json:"writeCount"`
	OtherCount uint64      `json:"otherCount"`
	ReadBytes  sysmon.Size `json:"readBytes"`
	WriteBytes sysmon.Size `json:"writeBytes"`
	OtherBytes sysmon.Size `json:"otherBytes"`
}

func (ioc IOCounters) GoString() string {
	s := []string{"IOCounters{", 
			fmt.Sprintf("  ReadCount  : %d", ioc.ReadCount), 
			fmt.Sprintf("  WriteCount : %d", ioc.WriteCount), 
			fmt.Sprintf("  OtherCount : %d", ioc.OtherCount), 
			fmt.Sprintf("  ReadBytes  : %s", ioc.ReadBytes), 
			fmt.Sprintf("  WriteBytes : %s", ioc.WriteBytes), 
			fmt.Sprintf("  OtherBytes : %s", ioc.OtherBytes), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func (p Process) ioCounters() (*IOCounters, error){
	wioc, err := win32.GetProcessIoCounters(syscall.Handle(p.handle)) 
	if err != nil {
		return nil, err
	}

	return &IOCounters{
		ReadCount  : wioc.ReadOperationCount,
		WriteCount : wioc.WriteOperationCount,
		OtherCount : wioc.OtherOperationCount,
		ReadBytes  : sysmon.Size(wioc.ReadTransferCount),
		WriteBytes : sysmon.Size(wioc.WriteTransferCount),
		OtherBytes : sysmon.Size(wioc.OtherTransferCount),
	}, nil
}

//---------------------------------------------------------------------------------------

type TimeUsage struct {
	CreationTime time.Time     `json:"creationTime"`
	ExitTime     time.Time     `json:"exitTime"`
	KernelTime   time.Duration `json:"kernelTime"`
	UserTime     time.Duration `json:"userTime"`
}

func (tu TimeUsage) GoString() string {
	s := []string{"TimeUsage{", 
			fmt.Sprintf("  CreationTime : %s", tu.CreationTime), 
			fmt.Sprintf("  ExitTime     : %s", tu.ExitTime), 
			fmt.Sprintf("  KernelTime   : %s", tu.KernelTime), 
			fmt.Sprintf("  UserTime     : %s", tu.UserTime),  
			"}",
	}
	return strings.Join(s, "\n")	
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

//---------------------------------------------------------------------------------------

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

//---------------------------------------------------------------------------------------

type Module struct {
	ProcessID uint32
	BaseAddr  uintptr      // The base address of the module in the context of the owning process.
	BaseSize  sysmon.Size  // The size of the module, in bytes.
	Handle    uintptr      // A handle to the module in the context of the owning process.
	Name      string
	ExePath   string
}

func (m Module) GoString() string {
	s := []string{"Module{", 
			fmt.Sprintf("  ProcessID : %v", m.ProcessID),   
			fmt.Sprintf("  BaseAddr  : %x", m.BaseAddr),   
			fmt.Sprintf("  BaseSize  : %s", m.BaseSize),   
			fmt.Sprintf("  Handle    : %v", m.Handle),   
			fmt.Sprintf("  Name      : %s", m.Name),   
			fmt.Sprintf("  ExePath   : %s", m.ExePath),   
			"}",
	}
	return strings.Join(s, "\n")	
}

func (p Process) modules() ([]*Module, error) {
	var ret []*Module

	snapshot, err := win32.CreateToolhelp32Snapshot(win32.TH32CS_SNAPMODULE, p.Pid)
	if err != nil {
		return ret, err
	}
	defer syscall.CloseHandle(snapshot)

	var modEntry win32.ModuleEntry32
	modEntry.Size = uint32(unsafe.Sizeof(modEntry))

	if err := win32.Module32First(snapshot, &modEntry); err != nil {
		return ret, err 
	}

	for {
		m := &Module{
			ProcessID : modEntry.ProcessID,
			BaseAddr  : modEntry.BaseAddr,
			BaseSize  : sysmon.Size(modEntry.BaseSize),
			Handle    : modEntry.Handle,
			Name      : syscall.UTF16ToString(modEntry.ModuleName[:]),
			ExePath   : syscall.UTF16ToString(modEntry.ExePath[:]),
		}
		ret = append(ret, m)

		err = win32.Module32Next(snapshot, &modEntry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return ret, err
		}
	}
	return ret[1:], nil
}

//---------------------------------------------------------------------------------------

type Thread struct {
	ThreadID        uint32
	OwnerProcessID  uint32
	BasePriority    int32  // The kernel base priority level assigned to the thread. 
	                       // The priority is a number from 0 to 31, with 0 representing 
	                       // the lowest possible thread priority.
}

func (t Thread) GoString() string {
	s := []string{"Thread{", 
			fmt.Sprintf("  ThreadID       : %v", t.ThreadID),   
			fmt.Sprintf("  OwnerProcessID : %v", t.OwnerProcessID),   
			fmt.Sprintf("  BasePriority   : %v", t.BasePriority),   
			"}",
	}
	return strings.Join(s, "\n")	
}

func (p Process) threads() ([]*Thread, error) {
	var ret []*Thread

	snapshot, err := win32.CreateToolhelp32Snapshot(win32.TH32CS_SNAPTHREAD, p.Pid)
	if err != nil {
		return ret, err
	}
	defer syscall.CloseHandle(snapshot)

	var thEntry win32.ThreadEntry32
	thEntry.Size = uint32(unsafe.Sizeof(thEntry))

	if err = win32.Thread32First(snapshot, &thEntry); err != nil {
		return ret, err
	}

	for {
		t := &Thread{
			ThreadID       : thEntry.ThreadID,
			OwnerProcessID : thEntry.OwnerProcessID,
			BasePriority   : thEntry.BasePriority,
		}
		ret = append(ret, t)

		err = win32.Thread32Next(snapshot, &thEntry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return ret, err
		}
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

func toDuration(ft syscall.Filetime) time.Duration {
	n := int64(ft.HighDateTime) << 32 + int64(ft.LowDateTime) // in 100-nanosecond intervals
	return time.Duration(n * 100) * time.Nanosecond
}
