
// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"
	"testing"
)

const PID = 10336

func TestFindProcess(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	fmt.Printf("%#v\n", p)
}

func TestUsage(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	u, err := p.Usage()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	fmt.Printf("%#v\n", u)
}

func TestIOCounters(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	ioc, err := p.IOCounters()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	fmt.Printf("%#v\n", ioc)
}

func TestMemoryInfo(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	pmc, err := p.MemoryInfo()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	fmt.Printf("%#v\n", pmc)
}

func TestModules(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	modules, err := p.Modules()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(modules) == 0 {
		t.Errorf("error: No modules!")
	}
/*
	for _, m := range modules {
		fmt.Printf("%#v\n", m)
	}
*/
}

func TestThreads(t *testing.T) {
	p, err := OpenProcess(PID)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	threads, err := p.Threads()
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if len(threads) == 0 {
		t.Errorf("error: No threads!")
	}
/*
	for _, t := range threads {
		fmt.Printf("%#v\n", t)
	}
*/
}
