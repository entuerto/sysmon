// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package disk

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/entuerto/sysmon"
	"github.com/entuerto/sysmon/internal/win32"
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

func allPartitions() ([]Partition, error) {
	var ret []Partition
	var dst []win32.Win32_DiskPartition

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
			Size            : sysmon.Size(p.Size),
			Index           : p.DiskIndex,
			Type            : p.Type,
		}

		ret = append(ret, part)
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

func allVolumes() ([]Volume, error) {
	var ret []Volume
	var dst []win32.Win32_Volume

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
			Capacity    : sysmon.Size(v.Capacity),
			Mount       : v.DriveLetter,
			DriveType   : v.DriveType,
			FileSystem  : v.FileSystem,
			FreeSpace   : sysmon.Size(v.FreeSpace),
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
	var dst []win32.Win32_Volume

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
			Total       : sysmon.Size(v.Capacity),
			Free        : sysmon.Size(v.FreeSpace),
			Used        : sysmon.Size(v.Capacity - v.FreeSpace),
			UsedPercent : float64(v.Capacity - v.FreeSpace) / float64(v.Capacity) * 100,
		}

		ret = append(ret, u)
	}

	return ret, nil
}

//---------------------------------------------------------------------------------------

func allDrives() ([]DiskDrive, error) {
	var ret []DiskDrive
	var dst []win32.Win32_DiskDrive

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
			Size           : sysmon.Size(d.Size),
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

func toDuration(v uint64) time.Duration {
	//DISK_PERFORMANCE::ReadTime is the total time spend on reads in 100ns units.
	return time.Duration(v * 100) * time.Nanosecond
}

func queryIOCounters(name string, freq time.Duration, qio QueryIO) {
	h, err := syscall.Open(name, syscall.O_RDONLY, 0)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <- time.After(freq):
			diskPerf, err := win32.IoctlDiskPerformance(h)

			if err != nil {
				log.Fatal(err)
			}

			ioc := &IOCounters{
				Name       :  name,
				ReadCount  :  diskPerf.ReadCount,
				WriteCount :  diskPerf.WriteCount,
				ReadBytes  :  sysmon.Size(diskPerf.BytesRead),
				WriteBytes :  sysmon.Size(diskPerf.BytesWritten),
				ReadTime   :  toDuration(diskPerf.ReadTime),
				WriteTime  :  toDuration(diskPerf.WriteTime),
				IoTime     :  toTime(diskPerf.QueryTime),
			}
			qio.IOCounterChan <- ioc
		case <- qio.quit:
			// we have received a signal to stop
			syscall.CloseHandle(h)
			return
		}
	}
}

func toTime(nsec int64) time.Time {
	// change starting time to the Epoch (00:00:00 UTC, January 1, 1970)
	nsec -= 116444736000000000
	// convert into nanoseconds
	nsec *= 100
	return time.Unix(0, nsec)
}
