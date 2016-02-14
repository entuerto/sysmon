// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procCreateToolhelp32Snapshot    = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procHeap32First                 = modkernel32.NewProc("Heap32FirstW") 
	procHeap32ListFirst             = modkernel32.NewProc("Heap32ListFirstW")
	procHeap32ListNext              = modkernel32.NewProc("Heap32ListNextW") 
	procHeap32Next                  = modkernel32.NewProc("Heap32NextW") 
	procModule32First               = modkernel32.NewProc("Module32FirstW") 
	procModule32Next                = modkernel32.NewProc("Module32NextW") 
	procProcess32First              = modkernel32.NewProc("Process32FirstW")
	procProcess32Next               = modkernel32.NewProc("Process32NextW")
	procThread32First               = modkernel32.NewProc("Thread32First") 
	procThread32Next                = modkernel32.NewProc("Thread32Next") 
	procToolhelp32ReadProcessMemory = modkernel32.NewProc("Toolhelp32ReadProcessMemory")

	procGetProcessHandleCount       = modkernel32.NewProc("GetProcessHandleCount")
	procGetProcessIoCounters        = modkernel32.NewProc("GetProcessIoCounters")
	procGetProcessMemoryInfo        = modkernel32.NewProc("K32GetProcessMemoryInfo")

)

const (
	MAX_PATH          = 260
	MAX_MODULE_NAME32 = 255
)

const (
	// flags for CreateToolhelp32Snapshot
	TH32CS_SNAPHEAPLIST = 0x01
	TH32CS_SNAPPROCESS  = 0x02
	TH32CS_SNAPTHREAD   = 0x04
	TH32CS_SNAPMODULE   = 0x08
	TH32CS_SNAPMODULE32 = 0x10
	TH32CS_SNAPALL      = TH32CS_SNAPHEAPLIST | TH32CS_SNAPMODULE | TH32CS_SNAPPROCESS | TH32CS_SNAPTHREAD
	TH32CS_INHERIT      = 0x80000000
)

type ModuleEntry32 struct {
	Size         uint32
	ModuleID     uint32  // This member is no longer used, and is always set to one.
	ProcessID    uint32
	GlblcntUsage uint32  // The load count of the module, which is not generally meaningful, and usually equal to 0xFFFF.
	ProccntUsage uint32  // The load count of the module, which is not generally meaningful, and usually equal to 0xFFFF.
	BaseAddr     uintptr // The base address of the module in the context of the owning process.
	BaseSize     uint32  // The size of the module, in bytes.
	Module       uintptr // A handle to the module in the context of the owning process.
	ModuleName   [MAX_MODULE_NAME32 + 1]uint16
	ExePath      [MAX_PATH]uint16
}

type ProcessEntry32 struct {
	Size            uint32
	Usage           uint32   // This member is no longer used and is always set to zero.
	ProcessID       uint32
	DefaultHeapID   uintptr  // This member is no longer used and is always set to zero.
	ModuleID        uint32   // This member is no longer used and is always set to zero.
	Threads         uint32   // The number of execution threads started by the process.
	ParentProcessID uint32
	PriClassBase    int32    // The base priority of any threads created by this process.
	Flags           uint32   // This member is no longer used and is always set to zero.
	ExeFile         [MAX_PATH]uint16 // To retrieve the full path to the executable file, 
	                                 // call the Module32First function and check the 
	                                 // szExePath member of the MODULEENTRY32 structure 
	                                 // that is returned.
}

type ThreadEntry32 struct {
	Size            uint32
	Usage           uint32  // This member is no longer used and is always set to zero.
	ThreadID        uint32
	OwnerProcessID  uint32
	BasePriority    int32  // The kernel base priority level assigned to the thread. 
	                       // The priority is a number from 0 to 31, with 0 representing 
	                       // the lowest possible thread priority.
	DeltaPriority   int32  // This member is no longer used and is always set to zero.
	Flags           uint32 // This member is no longer used and is always set to zero.
}

func CreateToolhelp32Snapshot(flags uint32, processId uint32) (handle syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procCreateToolhelp32Snapshot.Addr(), 2, uintptr(flags), uintptr(processId), 0)
	handle = syscall.Handle(r0)
	if handle == syscall.InvalidHandle {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Module32First(snapshot syscall.Handle, modEntry *ModuleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32First.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(modEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Module32Next(snapshot syscall.Handle, modEntry *ModuleEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procModule32Next.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(modEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Process32First(snapshot syscall.Handle, procEntry *ProcessEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procProcess32First.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(procEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Process32Next(snapshot syscall.Handle, procEntry *ProcessEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procProcess32Next.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(procEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Thread32First(snapshot syscall.Handle, thEntry *ThreadEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procThread32First.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(thEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func Thread32Next(snapshot syscall.Handle, thEntry *ThreadEntry32) (err error) {
	r1, _, e1 := syscall.Syscall(procThread32Next.Addr(), 2, uintptr(snapshot), uintptr(unsafe.Pointer(thEntry)), 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func FindProcessEntry(pid uint32) (*ProcessEntry32, error) {
	snapshot, err := CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	var procEntry ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))

	if err = Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}

	for {
		if procEntry.ProcessID == pid {
			return &procEntry, nil
		}
		err = Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

func FirstModuleEntry(pid uint32) (*ModuleEntry32, error) {
	snapshot, err := CreateToolhelp32Snapshot(TH32CS_SNAPMODULE, pid)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(snapshot)

	var modEntry ModuleEntry32
	modEntry.Size = uint32(unsafe.Sizeof(modEntry))

	if err := Module32First(snapshot, &modEntry); err != nil {
		return nil, err 
	}
	return &modEntry, nil
}

func GetProcessHandleCount(h syscall.Handle) (uint32, error) {
	var count uint32

	r, _, _ := procGetProcessHandleCount.Call(uintptr(h), uintptr(unsafe.Pointer(&count)))

	if r == 0 {
		return 0, syscall.GetLastError()
	}

	return count, nil
}

type IOCounters struct  {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

func GetProcessIoCounters(h syscall.Handle) (*IOCounters, error) {
	var ioc IOCounters

	r, _, _ := procGetProcessIoCounters.Call(uintptr(h), uintptr(unsafe.Pointer(&ioc)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &ioc, nil
}

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
