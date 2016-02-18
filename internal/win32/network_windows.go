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
	procGetExtendedUdpTable = modiphlpapi.NewProc("GetExtendedUdpTable")
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

type MIB_TCPTABLE_OWNER_PID struct {
	NumEntries uint32
	Table      []MIB_TCPROW_OWNER_PID
}

type MIB_TCPROW_OWNER_PID struct {
	State      uint32
	LocalAddr  [4]byte
	LocalPort  uint32
	RemoteAddr [4]byte
	RemotePort uint32
	OwningPid  uint32
}

// The GetExtendedTcpTableIPv4 function retrieves a table that contains a list of TCP 
// endpoints available to the application.
func GetExtendedTcpTableIPv4() (MIB_TCPTABLE_OWNER_PID, error) {
	var (
		TcpTable MIB_TCPTABLE_OWNER_PID
		TcpTableSize uint32
	)

	TcpTableSize = uint32(unsafe.Sizeof(TcpTable))

	tableSlice := make([]byte, TcpTableSize)

	r, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&tableSlice[0])),
		uintptr(unsafe.Pointer(&TcpTableSize)),
		uintptr(0),
		uintptr(syscall.AF_INET),
		uintptr(TCP_TABLE_OWNER_PID_ALL),
		uintptr(0))

	if syscall.Errno(r) == syscall.ERROR_INSUFFICIENT_BUFFER {
		tableSlice = make([]byte, TcpTableSize)

		r, _, _ = procGetExtendedTcpTable.Call(
			uintptr(unsafe.Pointer(&tableSlice[0])),
			uintptr(unsafe.Pointer(&TcpTableSize)),
			uintptr(0),
			uintptr(syscall.AF_INET),
			uintptr(TCP_TABLE_OWNER_PID_ALL),
			uintptr(0))
	}

	buf := bytes.NewReader(tableSlice)
	err := binary.Read(buf, binary.LittleEndian, &TcpTable.NumEntries)
	if err != nil {
		return TcpTable, err
	}

	TcpTable.Table = make([]MIB_TCPROW_OWNER_PID, TcpTable.NumEntries)

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

//----------------------------------------------------------------------------------

type MIB_TCP6TABLE_OWNER_PID struct {
	NumEntries uint32
	Table      []MIB_TCP6ROW_OWNER_PID
}

type MIB_TCP6ROW_OWNER_PID struct {
	LocalAddr     [16]byte
	LocalScopeId  uint32
	LocalPort     uint32
	RemoteAddr    [16]byte
	RemoteScopeId uint32
	RemotePort    uint32
	State         uint32
	OwningPid     uint32
}

// The GetExtendedTcpTableIPv6 function retrieves a table that contains a list of TCP 
// endpoints available to the application.
func GetExtendedTcpTableIPv6() (MIB_TCP6TABLE_OWNER_PID, error) {
	var (
		TcpTable MIB_TCP6TABLE_OWNER_PID
		TcpTableSize uint32
	)

	TcpTableSize = uint32(unsafe.Sizeof(TcpTable))

	tableSlice := make([]byte, TcpTableSize)

	r, _, _ := procGetExtendedTcpTable.Call(
		uintptr(unsafe.Pointer(&tableSlice[0])),
		uintptr(unsafe.Pointer(&TcpTableSize)),
		uintptr(0),
		uintptr(syscall.AF_INET6),
		uintptr(TCP_TABLE_OWNER_PID_ALL),
		uintptr(0))

	if syscall.Errno(r) == syscall.ERROR_INSUFFICIENT_BUFFER {
		tableSlice = make([]byte, TcpTableSize)

		r, _, _ = procGetExtendedTcpTable.Call(
			uintptr(unsafe.Pointer(&tableSlice[0])),
			uintptr(unsafe.Pointer(&TcpTableSize)),
			uintptr(0),
			uintptr(syscall.AF_INET6),
			uintptr(TCP_TABLE_OWNER_PID_ALL),
			uintptr(0))
	}

	buf := bytes.NewReader(tableSlice)
	err := binary.Read(buf, binary.LittleEndian, &TcpTable.NumEntries)
	if err != nil {
		return TcpTable, err
	}

	TcpTable.Table = make([]MIB_TCP6ROW_OWNER_PID, TcpTable.NumEntries)

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

//----------------------------------------------------------------------------------

const (
	UDP_TABLE_BASIC         = iota
	UDP_TABLE_OWNER_PID
	UDP_TABLE_OWNER_MODULE
)

type MIB_UDPTABLE_OWNER_PID struct {
	NumEntries uint32
	Table      []MIB_UDPROW_OWNER_PID
}

type MIB_UDPROW_OWNER_PID struct {
	LocalAddr  [4]byte
	LocalPort  uint32
	OwningPid  uint32
}

// The GetExtendedUdpTableIPv4 function retrieves a table that contains a list of UDP 
// endpoints available to the application.
func GetExtendedUdpTableIPv4() (MIB_UDPTABLE_OWNER_PID, error) {
		var (
		UdpTable MIB_UDPTABLE_OWNER_PID
		UdpTableSize uint32
	)

	UdpTableSize = uint32(unsafe.Sizeof(UdpTable))

	tableSlice := make([]byte, UdpTableSize)

	r, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&tableSlice[0])),
		uintptr(unsafe.Pointer(&UdpTableSize)),
		uintptr(0),
		uintptr(syscall.AF_INET),
		uintptr(UDP_TABLE_OWNER_PID),
		uintptr(0))

	if syscall.Errno(r) == syscall.ERROR_INSUFFICIENT_BUFFER {
		tableSlice = make([]byte, UdpTableSize)

		r, _, _ = procGetExtendedUdpTable.Call(
			uintptr(unsafe.Pointer(&tableSlice[0])),
			uintptr(unsafe.Pointer(&UdpTableSize)),
			uintptr(0),
			uintptr(syscall.AF_INET),
			uintptr(UDP_TABLE_OWNER_PID),
			uintptr(0))
	}

	buf := bytes.NewReader(tableSlice)
	err := binary.Read(buf, binary.LittleEndian, &UdpTable.NumEntries)
	if err != nil {
		return UdpTable, err
	}

	UdpTable.Table = make([]MIB_UDPROW_OWNER_PID, UdpTable.NumEntries)

	for i := 0; i < int(UdpTable.NumEntries); i++ {
		err := binary.Read(buf, binary.LittleEndian, &UdpTable.Table[i])
		if err != nil {
			break
		}
	}

	if r != 0 {
		return UdpTable, syscall.Errno(r)
	}
	return UdpTable, nil
}

//----------------------------------------------------------------------------------

type MIB_UDP6TABLE_OWNER_PID struct {
	NumEntries uint32
	Table      []MIB_UDP6ROW_OWNER_PID
}

type MIB_UDP6ROW_OWNER_PID struct {
	LocalAddr    [16]byte
	LocalScopeId uint32
	LocalPort    uint32
	OwningPid    uint32
}

// The GetExtendedUdpTableIPv6 function retrieves a table that contains a list of UDP 
// endpoints available to the application.
func GetExtendedUdpTableIPv6() (MIB_UDP6TABLE_OWNER_PID, error) {
		var (
		UdpTable MIB_UDP6TABLE_OWNER_PID
		UdpTableSize uint32
	)

	UdpTableSize = uint32(unsafe.Sizeof(UdpTable))

	tableSlice := make([]byte, UdpTableSize)

	r, _, _ := procGetExtendedUdpTable.Call(
		uintptr(unsafe.Pointer(&tableSlice[0])),
		uintptr(unsafe.Pointer(&UdpTableSize)),
		uintptr(0),
		uintptr(syscall.AF_INET),
		uintptr(UDP_TABLE_OWNER_PID),
		uintptr(0))

	if syscall.Errno(r) == syscall.ERROR_INSUFFICIENT_BUFFER {
		tableSlice = make([]byte, UdpTableSize)

		r, _, _ = procGetExtendedUdpTable.Call(
			uintptr(unsafe.Pointer(&tableSlice[0])),
			uintptr(unsafe.Pointer(&UdpTableSize)),
			uintptr(0),
			uintptr(syscall.AF_INET),
			uintptr(UDP_TABLE_OWNER_PID),
			uintptr(0))
	}

	buf := bytes.NewReader(tableSlice)
	err := binary.Read(buf, binary.LittleEndian, &UdpTable.NumEntries)
	if err != nil {
		return UdpTable, err
	}

	UdpTable.Table = make([]MIB_UDP6ROW_OWNER_PID, UdpTable.NumEntries)

	for i := 0; i < int(UdpTable.NumEntries); i++ {
		err := binary.Read(buf, binary.LittleEndian, &UdpTable.Table[i])
		if err != nil {
			break
		}
	}

	if r != 0 {
		return UdpTable, syscall.Errno(r)
	}
	return UdpTable, nil
}
