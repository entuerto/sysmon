// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cpu

import (
	"fmt"
	"log"
	"testing"
)

func TestCpuInfo(t *testing.T) {
	cpus, err := GetInfo()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(cpus) == 0 {
		t.Error("could not get CPU Info")
	}

	for _, cpu := range cpus {
		if cpu.ModelName == "" {
			t.Errorf("could not get CPU Info: %v", cpu)
		}
	}
}

func ExampleCpuInfo() {
	cpus, err := GetInfo()

	if err != nil {
		log.Fatal(err)
	}

	for _, cpu := range cpus {
		fmt.Println("Cpu:        ", cpu.CPU)
		fmt.Println("VendorId:   ", cpu.VendorId)
		fmt.Println("Family:     ", cpu.Family)
		fmt.Println("Model:      ", cpu.Model)
		fmt.Println("Stepping:   ", cpu.Stepping)
		fmt.Println("PhysicalId: ", cpu.PhysicalId)
		fmt.Println("CoreId:     ", cpu.CoreId)
		fmt.Println("Cores:      ", cpu.Cores)
		fmt.Println("ModelName:  ", cpu.ModelName)
		fmt.Println("Mhz:        ", cpu.Mhz)
		fmt.Println("CacheSize:  ", cpu.CacheSize)
		fmt.Println("Flags:      ", cpu.Flags)
	}
	// Output: _
}

func TestCpuUsagePercent(t *testing.T) {
	usage, err := UsagePercent()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(usage) == 0 {
		t.Error("No usage information")
	}
}

func ExampleCpuUsagePercent() {
	usage, err := UsagePercent()

	if err != nil {
		log.Fatal(err)
	}

	for _, u := range usage {
		fmt.Printf("%2.2f ", u)
	}
	// Output: _
}

func TestCpuSysTimes(t *testing.T) {
	times, err := SystemTimes()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if times == nil {
		t.Error("No Cpu time information")
	}
}

func ExampleCpuSysTimes() {
	t, err := SystemTimes()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("User:      ", t.User)
	fmt.Println("System:    ", t.System)
	fmt.Println("Idle:      ", t.Idle)
	fmt.Println("Nice:      ", t.Nice)
	fmt.Println("Iowait:    ", t.Iowait)
	fmt.Println("Irq:       ", t.Irq)
	fmt.Println("Softirq:   ", t.Softirq)
	fmt.Println("Steal:     ", t.Steal)
	fmt.Println("Guest:     ", t.Guest)
	fmt.Println("GuestNice: ", t.GuestNice)
	// Output: _
}
