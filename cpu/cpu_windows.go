// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"fmt"
	"syscall"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/entuerto/sysmon/internal/win32"
)

type Times struct {
	User      time.Duration `json:"user"`
	System    time.Duration `json:"system"`
	Idle      time.Duration `json:"idle"`
	Kernel    time.Duration `json:"kernel"`
}

func getInfo() ([]Info, error) {
	var ret []Info
	var dst []win32.Win32_Processor

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
		IdleTime   syscall.Filetime 
		KernelTime syscall.Filetime
		UserTime   syscall.Filetime
	)

	if err := win32.GetSystemTimes(&IdleTime, &KernelTime, &UserTime); err != nil {
		return nil, err
	}

	idle   := fileTimeToDuration(IdleTime)
	kernel := fileTimeToDuration(KernelTime)

	return &Times{
		Idle   : idle,
		User   : fileTimeToDuration(UserTime),
		System : kernel - idle,
		Kernel : kernel,
	}, nil
}

// A float representing the current system-wide CPU utilization as a percentage.
func usagePercent() ([]float64, error) {
	var ret []float64
	var dst []win32.Win32_Processor

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

func fileTimeToDuration(ft syscall.Filetime) time.Duration {
	n := int64(ft.HighDateTime) << 32 + int64(ft.LowDateTime) // in 100-nanosecond intervals
	return time.Duration(n * 100) * time.Nanosecond
}