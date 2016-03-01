// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"syscall"
	"unsafe"
)

const (
	IOCTL_DISK_PERFORMANCE = 0x70020
)

type DiskPerformance struct {
	BytesRead           uint64 // The number of bytes read.
	BytesWritten        uint64 // The number of bytes written.
	ReadTime            uint64 // The time it takes to complete a read.
	WriteTime           uint64 // The time it takes to complete a write.
	IdleTime            uint64 // The idle time.
	ReadCount           uint32 // The number of read operations.
	WriteCount          uint32 // The number of write operations.
	QueueDepth          uint32 // The depth of the queue.
	SplitCount          uint32 // The cumulative count of I/Os that are associated I/Os. 
	QueryTime            int64 // The system time stamp when a query for this structure 
	                           // is returned. 
	StorageDeviceNumber uint32 // The unique number for a device that identifies it to the 
	                           // storage manager that is indicated in the StorageManagerName 
	                           // member.
	StorageManagerName  [8]uint16 // The name of the storage manager that controls this device. 
}

func IoctlDiskPerformance(h syscall.Handle) (*DiskPerformance, error) {
	var (
		diskPerf DiskPerformance
		bytesReturned uint32
	)

	err := syscall.DeviceIoControl(h, 
				IOCTL_DISK_PERFORMANCE, 
				nil, 
				0, 
			    (*byte)(unsafe.Pointer(&diskPerf)), 
				uint32(unsafe.Sizeof(diskPerf)),
			    &bytesReturned, nil)

	if err != nil {
		return nil, err
	}

	return &diskPerf, nil
}
