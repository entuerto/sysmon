// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"fmt"
	"syscall"
	"time"

type Times struct {
	User      time.Duration `json:"user"`
	System    time.Duration `json:"system"`
	Idle      time.Duration `json:"idle"`
	Nice      time.Duration `json:"nice"`
	Iowait    time.Duration `json:"iowait"`
	Irq       time.Duration `json:"irq"`
	Softirq   time.Duration `json:"softirq"`
	Steal     time.Duration `json:"steal"`
	Guest     time.Duration `json:"guest"`
	GuestNice time.Duration `json:"guestNice"`
}

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
	cpuload, err := darwin.HostStatisticsCpuLoadInfo()
	if err != nil {
		return nil, err
	}

	user   := time.Duration(cpuload.cpu_ticks[darwin.CPU_STATE_USER] / darwin.CLK_TCK) * time.Second 
	system := time.Duration(cpuload.cpu_ticks[darwin.CPU_STATE_SYSTEM] / darwin.CLK_TCK) * time.Second
	idle   := time.Duration(cpuload.cpu_ticks[darwin.CPU_STATE_IDLE] / darwin.CLK_TCK) * time.Second
	nice   := time.Duration(cpuload.cpu_ticks[darwin.CPU_STATE_NICE] / darwin.CLK_TCK) * time.Second

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