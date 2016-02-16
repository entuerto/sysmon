// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"unsafe"
)

var (
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")

	procGetExtendedTcpTable = modiphlpapi.NewProc("GetExtendedTcpTable")
)

const (
	TCP_TABLE_BASIC_LISTENER           = iota
	TCP_TABLE_BASIC_CONNECTIONS        
	TCP_TABLE_BASIC_ALL                
	TCP_TABLE_OWNER_PID_LISTENER       
	TCP_TABLE_OWNER_PID_CONNECTIONS    
	TCP_TABLE_OWNER_PID_ALL            
	TCP_TABLE_OWNER_MODULE_LISTENER    
	TCP_TABLE_OWNER_MODULE_CONNECTIONS 
	TCP_TABLE_OWNER_MODULE_ALL         
)

type MIB_TCPTABLE struct {
	NumEntries uint32
    Table      []MIB_TCPROW
}

type MIB_TCPROW struct {
	State      uint32
	LocalAddr  [4]byte
	LocalPort  uint32
	RemoteAddr [4]byte
	RemotePort uint32
}

/*	
DWORD GetExtendedTcpTable(
  _Out_   PVOID           pTcpTable,
  _Inout_ PDWORD          pdwSize,
  _In_    BOOL            bOrder,     (FALSE)
  _In_    ULONG           ulAf,       (syscall.AF_INET)
  _In_    TCP_TABLE_CLASS TableClass, (TCP_TABLE_BASIC_ALL)
  _In_    ULONG           Reserved    (0)
);
*/

// The GetExtendedTcpTable function retrieves a table that contains a list of TCP 
// endpoints available to the application.
func GetExtendedTcpTable() (MIB_TCPTABLE, error) {
	var (
		TcpTable MIB_TCPTABLE
		TcpTableSize uint32
	)

	TcpTableSize = uint32(unsafe.Sizeof(TcpTable))

	tableSlice := make([]byte, TcpTableSize)

	r, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&tableSlice[0])),
		uintptr(unsafe.Pointer(&TcpTableSize)),
		uintptr(0),
		uintptr(syscall.AF_INET),
		uintptr(TCP_TABLE_BASIC_ALL),
		uintptr(0))

	if syscall.Errno(r) == syscall.ERROR_INSUFFICIENT_BUFFER {
		tableSlice = make([]byte, TcpTableSize)

		r, _, _ = procGetExtendedTcpTable.Call(
			uintptr(unsafe.Pointer(&tableSlice[0])),
			uintptr(unsafe.Pointer(&TcpTableSize)),
			uintptr(0),
			uintptr(syscall.AF_INET),
			uintptr(TCP_TABLE_BASIC_ALL),
			uintptr(0))
	}

	buf := bytes.NewReader(tableSlice)
	err := binary.Read(buf, binary.LittleEndian, &TcpTable.NumEntries)
	if err != nil {
		return TcpTable, err
	}

	TcpTable.Table = make([]MIB_TCPROW, TcpTable.NumEntries)

	for i := 0; i < int(TcpTable.NumEntries); i++ {
		err := binary.Read(buf, binary.LittleEndian, &TcpTable.Table[i])
		if err != nil {
			break
		}
	}

	if r != 0 {
		return TcpTable, syscall.Errno(r)
	}
	return TcpTable, nil
}