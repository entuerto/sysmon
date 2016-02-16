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


func queryConnections() ([]Connection, error) {
	table, err := win32.GetExtendedTcpTable() 
	if err != nil {
		return nil, err
	}

	var connections []Connection
	for _, r := range table.Table {
		c := Connection{
			State  : TcpState(r.State),
			Local  : &syscall.SockaddrInet4{
				Port : int(r.LocalPort),
				Addr : r.LocalAddr,
			},
			Remote : &syscall.SockaddrInet4{
				Port : int(r.RemotePort),
				Addr : r.RemoteAddr,
			},
		}
		connections = append(connections, c)
	}

	return connections, nil
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
