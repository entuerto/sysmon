// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

/*
#include <errno.h>
#include <stdbool.h>
#include <stdlib.h>
#include <stdio.h>
#include <utmpx.h>

#include <mach/mach.h>
#include <mach/task.h>
#include <mach/mach_init.h>
#include <mach/host_info.h>
#include <mach/mach_host.h>
#include <mach/mach_traps.h>
#include <mach/mach_vm.h>
#include <mach/shared_region.h>
*/
import "C"

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

const (
	CLK_TCK = 100  // ticks per second
)

func getInfo() ([]Info, error) {

	cpu, err := syscall.Sysctl("machdep.cpu.vendor")
	if err != nil {
		return nil, err
	}

	fmt.Println(cpu)
/*	
type Info struct {
	CPU        int32    `json:"cpu"`
	VendorId   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalId string   `json:"physicalId"`
	CoreId     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
}*/
	return nil, nil
}

func systemTimes() (*Times, error) {
	var (
		count C.mach_msg_type_number_t = C.HOST_CPU_LOAD_INFO_COUNT
		cpuload C.host_cpu_load_info_data_t
	)

	status := C.host_statistics(
		C.host_t(C.mach_host_self()),
		C.HOST_CPU_LOAD_INFO,
		C.host_info_t(unsafe.Pointer(&cpuload)),
		&count)

	if status != C.KERN_SUCCESS {
		return nil, fmt.Errorf("host_statistics error=%d", status)
	}

	user   := time.Duration(cpuload.cpu_ticks[C.CPU_STATE_USER] / CLK_TCK) * time.Second 
	system := time.Duration(cpuload.cpu_ticks[C.CPU_STATE_SYSTEM] / CLK_TCK) * time.Second
	idle   := time.Duration(cpuload.cpu_ticks[C.CPU_STATE_IDLE] / CLK_TCK) * time.Second
	nice   := time.Duration(cpuload.cpu_ticks[C.CPU_STATE_NICE] / CLK_TCK) * time.Second

    return &Times{
		Idle   : idle,
		User   : user,
		System : system,
		Nice   : nice,
	}, nil
    
}

func usagePercent() ([]float64, error) {
	return nil, nil
}