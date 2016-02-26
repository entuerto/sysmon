
// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/entuerto/sysmon/mem"
	"github.com/entuerto/sysmon/cpu"
)

/*
Field Description For Vm Mode

Procs
r: The number of processes waiting for run time.
b: The number of processes in uninterruptible sleep.

Memory
swpd:   the amount of virtual memory used.
free:   the amount of idle memory.
buff:   the amount of memory used as buffers.
cache:  the amount of memory used as cache.
inact:  the amount of inactive memory. 
active: the amount of active memory. 

Swap
si: Amount of memory swapped in from disk (/s).
so: Amount of memory swapped to disk (/s).

IO
bi: Blocks received from a block device (blocks/s).
bo: Blocks sent to a block device (blocks/s).

System
in: The number of interrupts per second, including the clock.
cs: The number of context switches per second.

CPU
These are percentages of total CPU time.
us: Time spent running non-kernel code. (user time, including nice time)
sy: Time spent running kernel code. (system time)
id: Time spent idle. Prior to Linux 2.5.41, this includes IO-wait time.
wa: Time spent waiting for IO. Prior to Linux 2.5.41, included in idle.
st: Time stolen from a virtual machine. Prior to Linux 2.6.11, unknown.
*/

/*
type Swap struct {
	Total   sysmon.Size `json:"total"` // total swap memory in bytes
	Used    sysmon.Size `json:"used"` // used swap memory in bytes
	Free    sysmon.Size `json:"free"` // free swap memory in bytes
	Percent float64     `json:"percent"` // the percentage usage calculated as (total - available) / total * 100
	SIn     sysmon.Size `json:"sin"` // the number of bytes the system has swapped in from disk (cumulative)
	SOut    sysmon.Size `json:"sout"` // the number of bytes the system has swapped out from disk (cumulative)
}
type Virtual struct {
	Total     sysmon.Size `json:"total"`
	Available sysmon.Size `json:"available"`
	Used      sysmon.Size `json:"used"`
	Free      sysmon.Size `json:"free"`
	Percent   float64     `json:"percent"` 
}
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
*/

type Data struct {
	swap *mem.Swap
	virt *mem.Virtual
	cpu  *cpu.Times
}

func CollectData(freq time.Duration) (chan *Data, chan bool) {
	DataChan := make(chan *Data)
	Quit     := make(chan bool)

	go func() {
		for {
			select {
			case <- time.After(freq):
				s, err := mem.SwapMemory() 
				if err != nil {
					log.Fatal(err)
				}

				v, err := mem.VirtualMemory() 
				if err != nil {
					log.Fatal(err)
				}

				c, err := cpu.SystemTimes() 
				if err != nil {
					log.Fatal(err)
				}

				DataChan <- &Data{
					swap : s,
					virt : v,
					cpu  : c,
				}
			case <- Quit:
				return
			}
		}
	}()

	return DataChan, Quit 
}

var (
	header1 = "procs -----------memory----------- ----swap---- -----io---- --system-- ------cpu------\n"
	header2 = " r  b   swpd   free   buff  cache    si     so    bi    bo    in   cs   us  sy  id  wa \n"
	line    = " 0  0 %6s %6s %6d %6d %5d %6d %5d %5d %5d %4d %4.0f%4.0f%4.0f%4.0f \n"
)



func main() {
	DataChan, QuitChan := CollectData(time.Second)

	fmt.Println()
	fmt.Printf(header1)
	fmt.Printf(header2)
	go func() {
		prevData := <- DataChan
		for {
			data := <- DataChan

			si := data.swap.SIn - prevData.swap.SIn
			so := data.swap.SOut - prevData.swap.SOut
			user := data.cpu.User - prevData.cpu.User
			sys  := data.cpu.System - prevData.cpu.System
			idle := data.cpu.Idle - prevData.cpu.Idle

			tot := user + sys + idle

			fmt.Printf(line, 
				       data.virt.Used, 
				       data.virt.Free, 
				       0, 
				       0, 
				       si, 
				       so,
				       0,
				       0,
				       0,
				       0,
				       float64(user) / float64(tot) * 100,
				       float64(sys) / float64(tot) * 100,
				       float64(idle) / float64(tot) * 100,
				       0.0)	

			prevData = data
		}
	}()

	<-time.After(20 * time.Second)
    QuitChan <- true
}
