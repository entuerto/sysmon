// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/entuerto/sysmon"
)

type Connections struct {

}

func QueryConnections() ([]Connections, error) {
	return QueryConnections()
}

//---------------------------------------------------------------------------------------

type IOCounters struct {
	Name        string      `json:"name"`
	BytesSent   sysmon.Size `json:"bytesSent"`   // number of bytes sent
	BytesRecv   sysmon.Size `json:"bytesRecv"`   // number of bytes received
	PacketsSent uint32      `json:"packetsSent"` // number of packets sent
	PacketsRecv uint32      `json:"packetsRecv"` // number of packets received
	Errin       uint32      `json:"errin"`        // total number of errors while receiving
	Errout      uint32      `json:"errout"`       // total number of errors while sending
	Dropin      uint32      `json:"dropin"`       // total number of incoming packets which were dropped
	Dropout     uint32      `json:"dropout"`      // total number of outgoing packets which were dropped (always 0 on OSX and BSD)
}

func (ioc IOCounters) GoString() string {
	s := []string{"IOCounters{", 
			fmt.Sprintf("  Name        : %s", ioc.Name), 
			fmt.Sprintf("  BytesSent   : %s", ioc.BytesSent), 
			fmt.Sprintf("  BytesRecv   : %s", ioc.BytesRecv), 
			fmt.Sprintf("  PacketsSent : %d", ioc.PacketsSent), 
			fmt.Sprintf("  PacketsRecv : %d", ioc.PacketsRecv), 
			fmt.Sprintf("  Errin       : %d", ioc.Errin), 
			fmt.Sprintf("  Errout      : %d", ioc.Errout), 
			fmt.Sprintf("  Dropin      : %d", ioc.Dropin), 
			fmt.Sprintf("  Dropout     : %d", ioc.Dropout), 
			"}",
	}
	return strings.Join(s, "\n")	
}

type QueryIO struct {
	IOCounterChan chan *IOCounters
	quit          chan bool
}

// Stop signals the goroutine to stop querying iocounters.
func (q QueryIO) Stop() {
	go func() {
		q.quit <- true
	}()
}

// diskperf.exe -Y and fisrt call to enable
func QueryIOCounters(iface net.Interface, freq time.Duration) (QueryIO, error) {
	qio := QueryIO{
		IOCounterChan : make(chan *IOCounters),
		quit          : make(chan bool),
	}

	go queryIOCounters(iface, freq, qio)

	return qio, nil 
}
