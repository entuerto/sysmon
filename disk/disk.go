// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package disk

import (
	"fmt"
	"strings"
	"time"

	"github.com/entuerto/sysmon"
)

type Device struct {
    DeviceId    string `json:"device_id"`
    
    Name        string `json:"name"`
    Caption     string `json:"caption"`
    Description string `json:"description"`
}

//---------------------------------------------------------------------------------------

type Partition struct {
	Device     
    
    BlockSize        uint64       `json:"blockSize"`
	BootPartition    bool         `json:"bootPartition"`
    NumberOfBlocks   uint64       `json:"numberOfBlocks"`
    PrimaryPartition bool         `json:"primaryPartition"`
    Size             sysmon.Size  `json:"size"`
    Index            uint32       `json:"index"`
    Status           string       `json:"status"`
    Type             string       `json:"type"`
}

func (p Partition) GoString() string {
	s := []string{"Partition{", 
			fmt.Sprintf("  DeviceId         : %s", p.DeviceId), 
			fmt.Sprintf("  Name             : %s", p.Name), 
			fmt.Sprintf("  Caption          : %s", p.Caption), 
			fmt.Sprintf("  Description      : %s", p.Description), 
			fmt.Sprintf("  BlockSize        : %d", p.BlockSize), 
			fmt.Sprintf("  BootPartition    : %t", p.BootPartition), 
			fmt.Sprintf("  NumberOfBlocks   : %d", p.NumberOfBlocks), 
			fmt.Sprintf("  PrimaryPartition : %t", p.PrimaryPartition),
			fmt.Sprintf("  Size             : %s", p.Size),  
			fmt.Sprintf("  Index            : %d", p.Index), 
			fmt.Sprintf("  Type             : %s", p.Type), 
			fmt.Sprintf("  Status           : %s", p.Status), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func AllPartitions() ([]Partition, error) {
	return allPartitions() 
}

//---------------------------------------------------------------------------------------

type DiskDrive struct {
	Device
    
    BytesPerSector uint32       `json:"bytesPerSector"`
    Partitions     uint32       `json:"partitions"`
    Model          string       `json:"model"`
    Size           sysmon.Size  `json:"size"`
    Index          uint32       `json:"index"`
    MediaType      string       `json:"mediaType"`
    SerialNumber   string       `json:"serialNumber"`
    Status         string       `json:"status"`
}

func (d DiskDrive) GoString() string {
	s := []string{"DiskDrive{", 
			fmt.Sprintf("  DeviceId       : %s", d.DeviceId), 
			fmt.Sprintf("  Name           : %s", d.Name), 
			fmt.Sprintf("  Caption        : %s", d.Caption), 
			fmt.Sprintf("  Description    : %s", d.Description), 
			fmt.Sprintf("  BytesPerSector : %d", d.BytesPerSector), 
			fmt.Sprintf("  Partitions     : %d", d.Partitions), 
			fmt.Sprintf("  Model          : %s", d.Model), 
			fmt.Sprintf("  Size           : %s", d.Size), 
			fmt.Sprintf("  Index          : %d", d.Index), 
			fmt.Sprintf("  MediaType      : %s", d.MediaType), 
			fmt.Sprintf("  SerialNumber   : %s", d.SerialNumber), 
			fmt.Sprintf("  Status         : %s", d.Status), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func AllDrives() ([]DiskDrive, error) {
	return allDrives() 
}

//---------------------------------------------------------------------------------------

type Volume struct {
	Device

	BlockSize       uint64      `json:"blockSize"`
	Capacity        sysmon.Size `json:"capacity"`
	DriveType       uint32      `json:"driveType"`
	FileSystem      string      `json:"fileSystem"`
	FreeSpace       sysmon.Size `json:"freeSpace"`
	Label           string      `json:"label"`
    Mount           string      `json:"mount"`
	SerialNumber    uint32      `json:"serialNumber"`
}

func (v Volume) GoString() string {
	s := []string{"Volume{", 
			fmt.Sprintf("  DeviceId     : %s", v.DeviceId), 
			fmt.Sprintf("  Name         : %s", v.Name), 
			fmt.Sprintf("  Caption      : %s", v.Caption), 
			fmt.Sprintf("  Description  : %s", v.Description), 
			fmt.Sprintf("  BlockSize    : %d", v.BlockSize), 
			fmt.Sprintf("  Capacity     : %s", v.Capacity), 
			fmt.Sprintf("  DriveType    : %d", v.DriveType), 
			fmt.Sprintf("  FileSystem   : %s", v.FileSystem), 
			fmt.Sprintf("  FreeSpace    : %s", v.FreeSpace), 
			fmt.Sprintf("  Label        : %s", v.Label), 
			fmt.Sprintf("  Mount        : %s", v.Mount), 
			fmt.Sprintf("  SerialNumber : %d", v.SerialNumber), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func AllVolumes() ([]Volume, error) {
	return allVolumes()
}

//---------------------------------------------------------------------------------------

type Usage struct {
	Device 
	Total        sysmon.Size  `json:"total"`
	Free         sysmon.Size  `json:"free"`
	Used         sysmon.Size  `json:"used"`
	UsedPercent  float64      `json:"usedPercent"`
}

func (u Usage) GoString() string {
	s := []string{"Usage{", 
			fmt.Sprintf("  DeviceId    : %s", u.DeviceId), 
			fmt.Sprintf("  Name        : %s", u.Name), 
			fmt.Sprintf("  Caption     : %s", u.Caption), 
			fmt.Sprintf("  Description : %s", u.Description), 
			fmt.Sprintf("  Total       : %s", u.Total), 
			fmt.Sprintf("  Free        : %s", u.Free), 
			fmt.Sprintf("  Used        : %s", u.Used), 
			fmt.Sprintf("  UsedPercent : %.2f", u.UsedPercent), 
			"}",
	}
	return strings.Join(s, "\n")	
}

func AllUsage() ([]Usage, error) {
	return allUsage()
}

//---------------------------------------------------------------------------------------

type IOCounters struct {
	Name         string      `json:"name"`
	ReadCount    uint32      `json:"readCount"`   // number of reads
	WriteCount   uint32      `json:"writeCount"`  // number of writes
	ReadBytes    sysmon.Size `json:"readBytes"`   // number of bytes read
	WriteBytes   sysmon.Size `json:"writeBytes"`  // number of bytes written
	ReadTime     uint64      `json:"readTime"`    // time spent reading from disk 
	WriteTime    uint64      `json:"writeTime"`   // time spent writing to disk 
	IoTime       uint64      `json:"ioTime"`
}

func (ioc IOCounters) GoString() string {
	s := []string{"IOCounters{", 
			fmt.Sprintf("  Name       : %s", ioc.Name), 
			fmt.Sprintf("  ReadCount  : %d", ioc.ReadCount), 
			fmt.Sprintf("  WriteCount : %d", ioc.WriteCount), 
			fmt.Sprintf("  ReadBytes  : %s", ioc.ReadBytes), 
			fmt.Sprintf("  WriteBytes : %s", ioc.WriteBytes), 
			fmt.Sprintf("  ReadTime   : %s", toDuration(ioc.ReadTime)), 
			fmt.Sprintf("  WriteTime  : %s", toDuration(ioc.WriteTime)), 
			fmt.Sprintf("  IoTime     : %d", ioc.IoTime), 
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
func QueryIOCounters(name string, freq time.Duration) (QueryIO, error) {
	qio := QueryIO{
		IOCounterChan : make(chan *IOCounters),
		quit          : make(chan bool),
	}

	go queryIOCounters(name, freq, qio)

	return qio, nil 
}
