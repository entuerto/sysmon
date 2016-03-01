// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"
	"strings"
)

// ProcessesByName()

type Process struct {
	Pid         uint32 `json:"pid"`
	ParentId    uint32 `json:"ppid"`
	Name        string `json:"name"`
	Executable  string `json:"executable"`
	CmdLine     string `json:"cmdLine"`
	HandleCount uint32 `json:"handleCount"`
	ThreadCount uint32 `json:"threadCount"`
	Status      string
	UserName    string

	handle uintptr  // windows proof
}

func (p Process) GoString() string {
	s := []string{"Process{", 
			fmt.Sprintf("  Pid         : %d", p.Pid),   
			fmt.Sprintf("  ParentId    : %d", p.ParentId),   
			fmt.Sprintf("  Name        : %s", p.Name),   
			fmt.Sprintf("  Executable  : %s", p.Executable),   
			fmt.Sprintf("  CmdLine     : %s", p.CmdLine),   
			fmt.Sprintf("  HandleCount : %d", p.HandleCount),   
			fmt.Sprintf("  ThreadCount : %d", p.ThreadCount),   
			fmt.Sprintf("  Status      : %s", p.Status),   
			fmt.Sprintf("  UserName    : %s", p.UserName),   
			"}",
	}
	return strings.Join(s, "\n")	
}

func (p Process) Parent() (*Process, error) {
	return OpenProcess(p.ParentId)
}

func (p Process) IOCounters() (*IOCounters, error) {
	return p.ioCounters()
}

func (p Process) MemoryInfo() (*MemoryCounters, error) {
	return p.memoryInfo()
}

func (p Process) Usage() (*TimeUsage, error) {
	return p.usage()
}

func (p Process) Modules() ([]*Module, error) {
	return p.modules()
}

func (p Process) Threads() ([]*Thread, error) {
	return p.threads()
}

// runtime.SetFinalizer(p, (*Process).Release)