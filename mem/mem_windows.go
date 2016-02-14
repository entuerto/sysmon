// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"syscall"
	"unsafe"

	"github.com/entuerto/sysmon"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGlobalMemoryStatusEx = modkernel32.NewProc("GlobalMemoryStatusEx")
)

type _MemoryStatusEx struct {
	Length               uint32
	MemoryLoad           uint32
	TotalPhys            uint64
	AvailPhys            uint64
	TotalPageFile        uint64
	AvailPageFile        uint64
	TotalVirtual         uint64
	AvailVirtual         uint64
	AvailExtendedVirtual uint64
}

func swapMemory() (*Swap, error) {
	var mem _MemoryStatusEx

	mem.Length = uint32(unsafe.Sizeof(mem))

	r, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &Swap{
		Total   : sysmon.Size(mem.TotalPageFile),
		Used    : sysmon.Size(mem.TotalPageFile - mem.AvailPageFile),
		Free    : sysmon.Size(mem.AvailPageFile), 
		Percent : float64(mem.TotalPageFile - mem.AvailPageFile) / float64(mem.TotalPageFile) * 100, // (total - available) / total * 100
	}, nil
}

func virtualMemory() (*Virtual, error) {
	var mem _MemoryStatusEx

	mem.Length = uint32(unsafe.Sizeof(mem))

	r, _, _ := procGlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&mem)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	return &Virtual{
		Total     : sysmon.Size(mem.TotalPhys),
		Available : sysmon.Size(mem.AvailPhys),
		Used      : sysmon.Size(mem.AvailPhys), 
		Free      : sysmon.Size(mem.TotalPhys - mem.AvailPhys), 
		Percent   : float64(mem.TotalPhys - mem.AvailPhys) / float64(mem.TotalPhys) * 100,
	}, nil
}