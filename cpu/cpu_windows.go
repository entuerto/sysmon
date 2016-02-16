// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
)

type Family uint16

func (f Family) String() string {
	switch (f) {
	case 1:
		return "Other"
	case 2:
		return "Unknown"
	case 198:
		return "Intel Core™ i7-2760QM"
	default:
		return fmt.Sprintf("%d", f)
	}
}

type Win32_Processor struct {
	Family                    Family
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	ProcessorId               *string
	MaxClockSpeed             uint32
}

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procGetSystemTimes = modkernel32.NewProc("GetSystemTimes")
)

func getInfo() ([]Info, error) {
	var ret []Info
	var dst []Win32_Processor

	q := wmi.CreateQuery(&dst, "")

	if err := wmi.Query(q, &dst); err != nil {
		return ret, err
	}

	fmt.Println("dst size: ", len(dst))

	var procID string
	for i, p := range dst {
		procID = ""
		if p.ProcessorId != nil {
			procID = *p.ProcessorId
		}

		cpu := Info{
			CPU:        int32(i),
			Family:     p.Family.String(),
			VendorId:   p.Manufacturer,
			ModelName:  p.Name,
			Cores:      int32(p.NumberOfLogicalProcessors),
			PhysicalId: procID,
			Mhz:        float64(p.MaxClockSpeed),
			Flags:      []string{},
		}
		ret = append(ret, cpu)
	}

	return ret, nil
}

// Retrieves system CPU timing information.
// On a multiprocessor system, the values returned are the sum of the designated 
// times across all processors.
func systemTimes() (*Times, error) {
	var (
		IdleTime syscall.Filetime 
		KernelTime syscall.Filetime
		UserTime syscall.Filetime
	)

	r, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(&IdleTime)),
		uintptr(unsafe.Pointer(&KernelTime)),
		uintptr(unsafe.Pointer(&UserTime)))

	if r == 0 {
		return nil, syscall.GetLastError()
	}

	idle   := FileTimeToDuration(IdleTime)
	kernel := FileTimeToDuration(KernelTime)

	return &Times{
		Idle:   idle,
		User:   FileTimeToDuration(UserTime),
		System: kernel - idle,
	}, nil
}

// A float representing the current system-wide CPU utilization as a percentage.
func usagePercent() ([]float64, error) {
	var ret []float64
	var dst []Win32_Processor

	q := wmi.CreateQuery(&dst, "")

	if err := wmi.Query(q, &dst); err != nil {
		return ret, err
	}
/*
	for _, p := range dst {
		// use range but windows can only get one percent.
		if p.LoadPercentage == nil {
			continue
		}
		ret = append(ret, float64(*p.LoadPercentage))
	}
*/	
	return ret, nil
}

func FileTimeToDuration(ft syscall.Filetime) time.Duration {
	n := int64(ft.HighDateTime) << 32 + int64(ft.LowDateTime) // in 100-nanosecond intervals
	return time.Duration(n * 100) * time.Nanosecond
}