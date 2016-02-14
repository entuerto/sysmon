// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"runtime"
	"time"
)

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
}

func Cores() int {
	return runtime.NumCPU()
}

func LogicalCores() int {
	return runtime.NumCPU()
}
