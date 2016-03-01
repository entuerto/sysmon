// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"fmt"
)

type Family uint16

func (f Family) String() string {
	switch (f) {
	case 1:
		return "Other"
	case 2:
		return "Unknown"
	case 198:
		return "Intel Core™ i7-2760QM"
	default:
		return fmt.Sprintf("%d", f)
	}
}

type Win32_Processor struct {
	Family                    Family
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	ProcessorId               *string
	MaxClockSpeed             uint32
}

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
	Size              uint64
	StartingOffset    uint64
	Type              string
}

type Win32_Volume  struct {
	Automount       bool
	BlockSize       uint64
	Capacity        uint64
	Caption         string
	Description     string
	DeviceID        string
	DriveLetter     string
	DriveType       uint32
	FileSystem      string
	FreeSpace       uint64
	Label           string
	Name            string
	NumberOfBlocks  uint64
	SystemName      string
	SerialNumber    uint32
}

type Win32_DiskDrive struct {
	DeviceID       string 
    Name           string
    Caption        string 
    Description    string 
    BytesPerSector uint32
    Partitions     uint32
    Model          string
    Size           uint64
    Index          uint32
    MediaType      string
    SerialNumber   string
    Status         string
}
