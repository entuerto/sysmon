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

	procGetSystemTimes       = modkernel32.NewProc("GetSystemTimes")
	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
)

type MemoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32 // A number between 0 and 100 that specifies the approximate percentage of physical memory that is in use
	TotalPhys            uint64 // The amount of actual physical memory, in bytes.
	AvailPhys            uint64 // The amount of physical memory currently available, in bytes.
	TotalPageFile        uint64 // The current committed memory limit for the system or the current process, whichever is smaller, in bytes.
	AvailPageFile        uint64 // The maximum amount of memory the current process can commit, in bytes.
	TotalVirtual         uint64 // The size of the user-mode portion of the virtual address space of the calling process, in bytes.
	AvailVirtual         uint64 // The amount of unreserved and uncommitted memory currently in the user-mode portion of the virtual address space of the calling process, in bytes.
	AvailExtendedVirtual uint64 // Reserved. This value is always 0.
}

func GetSystemTimes(IdleTime, KernelTime, UserTime *syscall.Filetime) error {
	r, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(IdleTime)),
		uintptr(unsafe.Pointer(KernelTime)),
		uintptr(unsafe.Pointer(UserTime)))

	if r == 0 {
		return syscall.GetLastError()
	}

	return nil
}

func GlobalMemoryStatusEx() (*MemoryStatusEx, error) {
	var mem MemoryStatusEx

	mem.Length = uint32(unsafe.Sizeof(mem))

	r, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &mem, nil
}