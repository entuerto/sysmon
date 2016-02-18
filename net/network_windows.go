// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/entuerto/sysmon"
	"github.com/entuerto/sysmon/internal/win32"
)

type ConnFilter func(c Connection) bool

func trueConnFilter(c Connection) bool {
	return true
}

func queryConnections(filter ConnFilter) ([]Connection, error) {
	var connections []Connection

	// IPv4
	tcpTable, err := win32.GetExtendedTcpTableIPv4() 
	if err != nil {
		return nil, err
	}

	
	for _, r := range tcpTable.Table {
		c := Connection{
			State  : ConnState(r.State),
			Type   : ConnTypeIPv4,
			Local  : &syscall.SockaddrInet4{
				Port : int(r.LocalPort),
				Addr : r.LocalAddr,
			},
			Remote : &syscall.SockaddrInet4{
				Port : int(r.RemotePort),
				Addr : r.RemoteAddr,
			},
			Pid    : r.OwningPid,
		}
		if filter(c) {
			connections = append(connections, c)
		}
	}

	// IPv6
	tcp6Table, err := win32.GetExtendedTcpTableIPv6() 
	if err != nil {
		return nil, err
	}

	for _, r := range tcp6Table.Table {
		c := Connection{
			State  : ConnState(r.State),
			Type   : ConnTypeIPv6,
			Local  : &syscall.SockaddrInet6{
				Port   : int(r.LocalPort),
				ZoneId : r.LocalScopeId,
				Addr   : r.LocalAddr,
			},
			Remote : &syscall.SockaddrInet6{
				Port   : int(r.RemotePort),
				ZoneId : r.RemoteScopeId,
				Addr   : r.RemoteAddr,
			},
			Pid    : r.OwningPid,
		}
		if filter(c) {
			connections = append(connections, c)
		}
	}

	// UDP4
	udpTable, err := win32.GetExtendedUdpTableIPv4() 
	if err != nil {
		return nil, err
	}

	for _, r := range udpTable.Table {
		c := Connection{
			Type   : ConnTypeUDP4,
			Local  : &syscall.SockaddrInet4{
				Port : int(r.LocalPort),
				Addr : r.LocalAddr,
			},
			Pid    : r.OwningPid,
		}
		if filter(c) {
			connections = append(connections, c)
		}
	}

	// UDP6
	udp6Table, err := win32.GetExtendedUdpTableIPv6() 
	if err != nil {
		return nil, err
	}

	for _, r := range udp6Table.Table {
		c := Connection{
			Type   : ConnTypeUDP6,
			Local  : &syscall.SockaddrInet6{
				Port   : int(r.LocalPort),
				ZoneId : r.LocalScopeId,
				Addr   : r.LocalAddr,
			},
			Pid    : r.OwningPid,
		}
		if filter(c) {
			connections = append(connections, c)
		}
	}

	return connections, nil
}

func queryConnectionsAll() ([]Connection, error) {
	return queryConnections(trueConnFilter)
}

func queryConnectionsByPid(Pid uint32) ([]Connection, error) {
	pidConnFilter := func (c Connection) bool {
		return c.Pid == Pid
	}
	return queryConnections(pidConnFilter) 
}

func queryIOCounters(iface net.Interface, freq time.Duration, qio QueryIO) {

	for {
		var row syscall.MibIfRow

		select {
		case <- time.After(freq):
			row.Index = uint32(iface.Index)

			err := syscall.GetIfEntry(&row)
			if err != nil {
				log.Fatal(err)
			}

			ioc := &IOCounters{
				Name        :  iface.Name,
				BytesSent   :  sysmon.Size(row.OutOctets),
				BytesRecv   :  sysmon.Size(row.InOctets),
				PacketsSent :  row.OutUcastPkts,
				PacketsRecv :  row.InUcastPkts,
				Errin       :  row.InErrors,
				Errout      :  row.OutErrors,
				Dropin      :  row.InDiscards,
				Dropout     :  row.OutDiscards,
			}
			qio.IOCounterChan <- ioc
		case <- qio.quit:
			// we have received a signal to stop
			return
		}
	}
}
