// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package disk

import (
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/StackExchange/wmi"
	"github.com/entuerto/sysmon"
)

type Access uint16

func (a Access) String() string {
	switch (a) {
	case 0:
		return "Unknown"
	case 1:
		return "Readable"
	case 2:
		return "Writable"
	case 3:
		return "Read/write"
	case 4:
		return "Write Once"
	default:
		return fmt.Sprintf("%d", a)
	}
}

//---------------------------------------------------------------------------------------

type Win32_DiskPartition struct {
	BlockSize         uint64
	BootPartition     bool
	PrimaryPartition  bool
	Caption           string
	CreationClassName string
	Description       string
	DeviceID          string
	DiskIndex         uint32
	Name              string
	NumberOfBlocks    uint64
	Size              sysmon.Size
	StartingOffset    uint64
	Type              string
}

func allPartitions() ([]Partition, error) {
	var ret []Partition
	var dst []Win32_DiskPartition

	q := wmi.CreateQuery(&dst, "")

	if err := wmi.Query(q, &dst); err != nil {
		return ret, err
	}

	for _, p := range dst {

		part := Partition {
			Device: Device{
				DeviceId    : p.DeviceID,
				Name        : p.Name,
				Caption     : p.Caption,
				Description : p.Description,
			},
			BlockSize       : p.BlockSize,
			BootPartition   : p.BootPartition,
			NumberOfBlocks  : p.NumberOfBlocks,
			PrimaryPartition: p.PrimaryPartition,
			Size            : p.Size,
			Index           : p.DiskIndex,
			Type            : p.Type,
		}

		ret = append(ret, part)
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

type Win32_Volume  struct {
	Automount       bool
	BlockSize       uint64
	Capacity        sysmon.Size
	Caption         string
	Description     string
	DeviceID        string
	DriveLetter     string
	DriveType       uint32
	FileSystem      string
	FreeSpace       sysmon.Size
	Label           string
	Name            string
	NumberOfBlocks  uint64
	SystemName      string
	SerialNumber    uint32
}

func allVolumes() ([]Volume, error) {
	var ret []Volume
	var dst []Win32_Volume

	q := wmi.CreateQuery(&dst, "")

	c := &wmi.Client{
		NonePtrZero: true,
		PtrNil: false,
		AllowMissingFields: true,
	}

	if err := c.Query(q, &dst); err != nil {
		return ret, err
	}

	for _, v := range dst {

		vol := Volume {
			Device: Device{
				DeviceId    : v.DeviceID,
				Name        : v.Name,
				Caption     : v.Caption,
				Description : v.Description,
			},
			BlockSize   : v.BlockSize,
			Capacity    : v.Capacity,
			Mount       : v.DriveLetter,
			DriveType   : v.DriveType,
			FileSystem  : v.FileSystem,
			FreeSpace   : v.FreeSpace,
			Label       : v.Label,
			SerialNumber: v.SerialNumber,   
		}

		ret = append(ret, vol)
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

func allUsage() ([]Usage, error) {
	var ret []Usage
	var dst []Win32_Volume

	q := wmi.CreateQuery(&dst, "")

	c := &wmi.Client{
		NonePtrZero: true,
		PtrNil: false,
		AllowMissingFields: true,
	}

	if err := c.Query(q, &dst); err != nil {
		return ret, err
	}

	for _, v := range dst {

		u := Usage {
			Device: Device{
				DeviceId    : v.DeviceID,
				Name        : v.Name,
				Caption     : v.Caption,
				Description : v.Description,
			},
			Total       : v.Capacity,
			Free        : v.FreeSpace,
			Used        : v.Capacity - v.FreeSpace,
			UsedPercent : float64(v.Capacity - v.FreeSpace) / float64(v.Capacity) * 100,
		}

		ret = append(ret, u)
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

type Win32_DiskDrive struct {
	DeviceID       string 
    Name           string
    Caption        string 
    Description    string 
    BytesPerSector uint32
    Partitions     uint32
    Model          string
    Size           sysmon.Size
    Index          uint32
    MediaType      string
    SerialNumber   string
    Status         string
}

func allDrives() ([]DiskDrive, error) {
	var ret []DiskDrive
	var dst []Win32_DiskDrive

	q := wmi.CreateQuery(&dst, "")

	c := &wmi.Client{
		NonePtrZero: true,
		PtrNil: false,
		AllowMissingFields: true,
	}

	if err := c.Query(q, &dst); err != nil {
		return ret, err
	}

	for _, d := range dst {

		disk := DiskDrive {
			Device: Device{
				DeviceId    : d.DeviceID,
				Name        : d.Name,
				Caption     : d.Caption,
				Description : d.Description,
			},
			BytesPerSector : d.BytesPerSector,
			Partitions     : d.Partitions,
			Model          : d.Model,
			Size           : d.Size,
			Index          : d.Index,
			MediaType      : d.MediaType,
			SerialNumber   : d.SerialNumber,
			Status         : d.Status,   
		}

		ret = append(ret, disk)
	}

	return ret, nil
} 

//---------------------------------------------------------------------------------------

type _DISK_PERFORMANCE struct {
	BytesRead           sysmon.Size
	BytesWritten        sysmon.Size
	ReadTime            uint64
	WriteTime           uint64
	IdleTime            uint64
	ReadCount           uint32
	WriteCount          uint32
	QueueDepth          uint32
	SplitCount          uint32
	QueryTime           uint64
	StorageDeviceNumber uint32
	StorageManagerName  [8]uint16
}

func (dp _DISK_PERFORMANCE) GoString() string {
	s := []string{"_DISK_PERFORMANCE{", 
			fmt.Sprintf("  BytesRead           : %s", dp.BytesRead), 
			fmt.Sprintf("  BytesWritten        : %s", dp.BytesWritten), 
			fmt.Sprintf("  ReadTime            : %s", toDuration(dp.ReadTime)), 
			fmt.Sprintf("  WriteTime           : %s", toDuration(dp.WriteTime)), 
			fmt.Sprintf("  IdleTime            : %s", toDuration(dp.IdleTime)), 
			fmt.Sprintf("  ReadCount           : %d", dp.ReadCount), 
			fmt.Sprintf("  QueueDepth          : %d", dp.QueueDepth), 
			fmt.Sprintf("  SplitCount          : %d", dp.SplitCount), 
			fmt.Sprintf("  QueryTime           : %s", toDuration(dp.QueryTime)), 
			fmt.Sprintf("  StorageDeviceNumber : %d", dp.StorageDeviceNumber), 
			fmt.Sprintf("  BytesRead           : %d", dp.BytesRead), 
			fmt.Sprintf("  StorageManagerName  : %s", syscall.UTF16ToString(dp.StorageManagerName[:])), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func toDuration(v uint64) time.Duration {
	return time.Duration(v * 100) * time.Nanosecond
}

//DISK_PERFORMANCE::ReadTime is the total time spend on reads in 100ns units.

const _IOCTL_DISK_PERFORMANCE = 0x70020

func queryIOCounters(name string, freq time.Duration, qio QueryIO) {
	h, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	var (
		diskPerf _DISK_PERFORMANCE
		bytesReturned uint32
	)

	for {
		select {
		case <- time.After(freq):
			err := syscall.DeviceIoControl(h, 
				_IOCTL_DISK_PERFORMANCE, 
				nil, 
				0, 
			    (*byte)(unsafe.Pointer(&diskPerf)), 
				uint32(unsafe.Sizeof(diskPerf)),
			    &bytesReturned, nil)
		
			if err != nil {
				log.Fatal(err)
			}

			ioc := &IOCounters{
				Name       :  name,
				ReadCount  :  diskPerf.ReadCount,
				WriteCount :  diskPerf.WriteCount,
				ReadBytes  :  diskPerf.BytesRead,
				WriteBytes :  diskPerf.BytesWritten,
				ReadTime   :  diskPerf.ReadTime,
				WriteTime  :  diskPerf.WriteTime,
				IoTime     :  0,
			}
			qio.IOCounterChan <- ioc
		case <- qio.quit:
			// we have received a signal to stop
			syscall.CloseHandle(h)
			return
		}
	}
}