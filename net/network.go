// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/entuerto/sysmon"
)

// Defines the connection state
type ConnState uint32

const (
	ConnStateClosed    ConnState = iota + 1
	ConnStateListen       
	ConnStateSynSent      
	ConnStateSynRcvd      
	ConnStateEstab        
	ConnStateFinWait1     
	ConnStateFinWait2     
	ConnStateCloseWait    
	ConnStateClosing      
	ConnStateLastAck      
	ConnStateTimeWait     
	ConnStateDeleteTcb    
)

func (s ConnState) String() string {
	switch s {
	case ConnStateClosed:
		return "CLOSED"
	case ConnStateListen:
		return "LISTEN"    
	case ConnStateSynSent:
		return "SYN-SENT"
	case ConnStateSynRcvd:
		return "SYN-RECEIVED"
	case ConnStateEstab:
		return "ESTABLISHED"
	case ConnStateFinWait1:
		return "FIN-WAIT-1"
	case ConnStateFinWait2:
		return "FIN-WAIT-2"
	case ConnStateCloseWait:
		return "CLOSE-WAIT"   
	case ConnStateClosing:
		return "CLOSING"     
	case ConnStateLastAck:
		return "LAST-ACK"
	case ConnStateTimeWait:
		return "TIME-WAIT"    
	case ConnStateDeleteTcb:
		return "DELETE-TCB"
	}
	return "UNKNOWN"
}

type ConnType uint32

const (
	ConnTypeIPv4 = iota + 1
	ConnTypeIPv6
	ConnTypeUDP4
	ConnTypeUDP6
	ConnTypeUnix
)

func (t ConnType) String() string {
	switch t {
	case ConnTypeIPv4:
		return "IPv4"
	case ConnTypeIPv6:
		return "IPv6"
	case ConnTypeUDP4:
		return "UDP4"
	case ConnTypeUDP6:
		return "UDP6"
	case ConnTypeUnix:
		return "Unix"
	}
	return "UNKNOWN"
}

type Connection struct {
	State   ConnState
	Type    ConnType
	Local   syscall.Sockaddr
	Remote  syscall.Sockaddr
	Pid     uint32
}

func QueryConnectionsAll() ([]Connection, error) {
	return queryConnectionsAll()
}

func QueryConnectionsByPid(Pid uint32) ([]Connection, error) {
	return queryConnectionsByPid(Pid)
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
