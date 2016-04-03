// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package darwin

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
	"unsafe"
)

const (
	CLK_TCK = 100  // ticks per second

	CPU_STATE_USER    = 0
	CPU_STATE_SYSTEM  = 1
	CPU_STATE_IDLE    = 2
	CPU_STATE_NICE    = 3
)

type HostCpuLoadInfoData C.host_cpu_load_info_data_t

func HostStatisticsCpuLoadInfo() (*HostCpuLoadInfoData, error) {
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

    return (*HostCpuLoadInfoData)(&cpuload), nil   
}