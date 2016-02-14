// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"log"
	"net"
	"syscall"
	"time"
	"unsafe"

	"github.com/entuerto/sysmon"
)

var (
	modiphlpapi = NewLazyDLL("iphlpapi.dll")

	procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")
)

type _MIB_TCPTABLE struct {
	NumEntries uint32
    Table      [ANY_SIZE]_MIB_TCPROW
}

func queryConnections() ([]Connections, error) {
/*	
DWORD GetExtendedTcpTable(
  _Out_   PVOID           pTcpTable,
  _Inout_ PDWORD          pdwSize,
  _In_    BOOL            bOrder,
  _In_    ULONG           ulAf,
  _In_    TCP_TABLE_CLASS TableClass,
  _In_    ULONG           Reserved
);
*/
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
