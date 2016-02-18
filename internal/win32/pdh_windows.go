// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"errors"
	"syscall"
	"time"
	"unsafe"
)

const (
	PDH_FMT_RAW      = 0x00000010
	PDH_FMT_ANSI     = 0x00000020
	PDH_FMT_UNICODE  = 0x00000040
	PDH_FMT_LONG     = 0x00000100
	PDH_FMT_DOUBLE   = 0x00000200
	PDH_FMT_LARGE    = 0x00000400
	PDH_FMT_NOSCALE  = 0x00001000
	PDH_FMT_1000     = 0x00002000
	PDH_FMT_NODATA   = 0x00004000
	PDH_FMT_NOCAP100 = 0x00008000

	PERF_DETAIL_NOVICE   = 100
	PERF_DETAIL_ADVANCED = 200
	PERF_DETAIL_EXPERT   = 300
	PERF_DETAIL_WIZARD   = 400
)

var (
	modpdh      = syscall.NewLazyDLL("pdh.dll")

	procPdhOpenQuery                = modpdh.NewProc("PdhOpenQuery")
	procPdhAddCounter               = modpdh.NewProc("PdhAddCounterW")
	procPdhAddEnglishCounter        = modpdh.NewProc("PdhAddEnglishCounterW")
	procPdhCollectQueryData         = modpdh.NewProc("PdhCollectQueryData")
	procPdhCollectQueryDataWithTime = modpdh.NewProc("PdhCollectQueryDataWithTime")
	procPdhEnumObjects              = modpdh.NewProc("PdhEnumObjectsW")
	procPdhEnumObjectItems          = modpdh.NewProc("PdhEnumObjectItemsW")
	procPdhGetFormattedCounterValue = modpdh.NewProc("PdhGetFormattedCounterValue")
	procPdhCloseQuery               = modpdh.NewProc("PdhCloseQuery")
	procPdhLookupPerfNameByIndex    = modpdh.NewProc("PdhLookupPerfNameByIndexW")
	procPdhLookupPerfIndexByName    = modpdh.NewProc("PdhLookupPerfIndexByNameW")
	procPdhRemoveCounter            = modpdh.NewProc("PdhRemoveCounter")
)

type PdhFmtCounterValueDouble struct {
	Status uint32
	Value  float64
}

type PdhFmtCounterValueLarge struct {
	Status uint32
	Value  int64
}

type PdhFmtCounterValueLong struct {
	Status  uint32
	Value   int32
	filler  int32
}

type PdhObjectItems struct {
	ObjectName string
	Counters   []string
	Instances  []string
}

// Creates a new query that is used to manage the collection of performance data.
func PdhOpenQuery() (syscall.Handle, error) {
	var h syscall.Handle

	r, _, err := procPdhOpenQuery.Call(0, 0, uintptr(unsafe.Pointer(&h)))
	if r != 0 {
		return 0, err
	}

	return h, nil
}

// Closes all counters contained in the specified query, closes all handles related to 
// the query, and frees all memory associated with the query.
func PdhCloseQuery(queryHdl syscall.Handle) error {
	r, _, err := procPdhAddCounter.Call(uintptr(queryHdl))
	if r != 0 {
		return err
	}
	return nil
}

// Adds the specified counter to the query. 
func PdhAddCounter(queryHdl syscall.Handle, CounterPath string) (syscall.Handle, error) {
	var counterHdl syscall.Handle

	r, _, err := procPdhAddCounter.Call(
					uintptr(queryHdl),
					uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(CounterPath))),
					0,
					uintptr(unsafe.Pointer(&counterHdl)))

	if r != 0 {
		return 0, err
	}

	return counterHdl, nil
}

// Adds the specified language-neutral counter to the query.
// See for more information: https://msdn.microsoft.com/en-us/library/windows/desktop/aa372536(v=vs.85).aspx
func PdhAddEnglishCounter(queryHdl syscall.Handle, CounterPath string) (syscall.Handle, error) {
	var counterHdl syscall.Handle

	r, _, err := procPdhAddEnglishCounter.Call(
					uintptr(queryHdl),
					uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(CounterPath))),
					0,
					uintptr(unsafe.Pointer(&counterHdl)))

	if r != 0 {
		return 0, err
	}

	return counterHdl, nil
}

func PdhGetCounterValueFloat64(counterHdl syscall.Handle) (float64, error) {
	var value PdhFmtCounterValueDouble

	r, _, _ := procPdhGetFormattedCounterValue.Call(
					uintptr(counterHdl), 
					PDH_FMT_DOUBLE, 
					uintptr(0), 
					uintptr(unsafe.Pointer(&value)))

	if r != 0 && r != uintptr(PDH_INVALID_DATA) {
		return 0.0, nil
	}

	return float64(value.Value), nil
}

func PdhGetCounterValueInt64(counterHdl syscall.Handle) (int64, error) {
	var value PdhFmtCounterValueLarge

	r, _, _ := procPdhGetFormattedCounterValue.Call(
					uintptr(counterHdl), 
					PDH_FMT_LARGE, 
					uintptr(0), 
					uintptr(unsafe.Pointer(&value)))

	if r != 0 && r != uintptr(PDH_INVALID_DATA) {
		return 0.0, nil
	}

	return int64(value.Value), nil
}

func PdhGetCounterValueInt32(counterHdl syscall.Handle) (int32, error) {
	var value PdhFmtCounterValueLong

	r, _, _ := procPdhGetFormattedCounterValue.Call(
					uintptr(counterHdl), 
					PDH_FMT_LONG, 
					uintptr(0), 
					uintptr(unsafe.Pointer(&value)))

	if r != 0 && r != uintptr(PDH_INVALID_DATA) {
		return 0.0, nil
	}

	return int32(value.Value), nil
}

// Collects the current raw data value for all counters in the specified query and 
// updates the status code of each counter. 
func PdhCollectQueryData(queryHdl syscall.Handle) error {
	r, _, err := procPdhCollectQueryData.Call(uintptr(queryHdl))

	if r != 0 && err != nil {
		if r == uintptr(PDH_NO_DATA) {
			// Has not values
			return errors.New(codeText[PDH_NO_DATA])
		}
		return err
	}

	return nil
}

// Collects the current raw data value for all counters in the specified query and 
// updates the status code of each counter.  
//
// Time stamp when the first counter value in the query was retrieved. The time is specified as FILETIME.
func PdhCollectQueryDataWithTime(queryHdl syscall.Handle) (*time.Time, error) {
	var TimeStamp syscall.Filetime

	r, _, err := procPdhCollectQueryDataWithTime.Call(
		uintptr(queryHdl),
		uintptr(unsafe.Pointer(&TimeStamp)))

	if r != 0 && err != nil {
		if r == uintptr(PDH_NO_DATA) {
			// Has not values
			return nil, errors.New(codeText[PDH_NO_DATA])
		}
		return nil, err
	}

	ts := time.Unix(0, TimeStamp.Nanoseconds())

	return &ts, nil
}

// Returns a list of objects available on the specified computer or in the specified 
// log file. 
func PdhEnumObjects(DataSource, MachineName string) ([]string, error) {
	var (
		ObjectList   []uint16
		BufferLength = uint32(0)

		source  = uintptr(0)
		machine = uintptr(0)
	)

	if len(DataSource) != 0 {
		ds, _ := syscall.UTF16PtrFromString(DataSource)
		source = uintptr(unsafe.Pointer(ds))
	}

	if len(MachineName) != 0 {
		m, _ := syscall.UTF16PtrFromString(MachineName)
		machine = uintptr(unsafe.Pointer(m))
	}

	r, _, _ := procPdhEnumObjects.Call(
		source, // real-time source,
		machine, // local machine
		uintptr(0),
		uintptr(unsafe.Pointer(&BufferLength)),
		uintptr(PERF_DETAIL_WIZARD),
		0)

	if r != PDH_MORE_DATA {
		return nil, errors.New(codeText[r])
	}

	ObjectList = make([]uint16, BufferLength)

	r, _, _ = procPdhEnumObjects.Call(
		source, // real-time source,
		machine, // local machine
		uintptr(unsafe.Pointer(&ObjectList[0])), // NULL to get size
		uintptr(unsafe.Pointer(&BufferLength)),
		uintptr(PERF_DETAIL_WIZARD),
		0)

	if r != 0 {
		return nil, errors.New(codeText[r])
	}

	var ret []string

	start := 0
	for i := 0; i < int(BufferLength); i++ {
		if ObjectList[i] == 0 {
			ret = append(ret, syscall.UTF16ToString(ObjectList[start:i]))
			start = i + 1
		}
	}
	return ret, nil
}

// Returns the specified object's counter and instance names that exist on the 
// specified computer or in the specified log file. 
func PdhEnumObjectItems(DataSource, MachineName, ObjectName string) (*PdhObjectItems, error) {
	var (
		source  = uintptr(0)
		machine = uintptr(0)

		CounterList       []uint16
		CounterListLength = uint32(0)

		InstanceList       []uint16
		InstanceListLength = uint32(0)
	)

	if len(DataSource) != 0 {
		ds, _ := syscall.UTF16PtrFromString(DataSource)
		source = uintptr(unsafe.Pointer(ds))
	}

	if len(MachineName) != 0 {
		m, _ := syscall.UTF16PtrFromString(MachineName)
		machine = uintptr(unsafe.Pointer(m))
	}

	r, _, _ := procPdhEnumObjectItems.Call(
		source, // real-time source,
		machine, // local machine
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(ObjectName))),
		uintptr(0),
		uintptr(unsafe.Pointer(&CounterListLength)),
		uintptr(0),
		uintptr(unsafe.Pointer(&InstanceListLength)),
		uintptr(PERF_DETAIL_WIZARD),
		0)

	if r != PDH_MORE_DATA {
		return nil, errors.New(codeText[r])
	}

	CounterList  = make([]uint16, CounterListLength)
	InstanceList = make([]uint16, InstanceListLength)

	r, _, _ = procPdhEnumObjectItems.Call(
		source, // real-time source,
		machine, // local machine
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(ObjectName))),
		uintptr(unsafe.Pointer(&CounterList[0])),
		uintptr(unsafe.Pointer(&CounterListLength)),
		uintptr(unsafe.Pointer(&InstanceList[0])),
		uintptr(unsafe.Pointer(&InstanceListLength)),
		uintptr(PERF_DETAIL_WIZARD),
		0)

	if r != 0 {
		return nil, errors.New(codeText[r])
	}

	var counters []string
	var instances []string

	start := 0
	for i := 0; i < int(CounterListLength); i++ {
		if CounterList[i] == 0 {
			counters = append(counters, syscall.UTF16ToString(CounterList[start:i]))
			start = i + 1
		}
	}

	start = 0
	for i := 0; i < int(InstanceListLength); i++ {
		if InstanceList[i] == 0 {
			instances = append(instances, syscall.UTF16ToString(InstanceList[start:i]))
			start = i + 1
		}
	}

	return &PdhObjectItems{
		ObjectName : ObjectName,
		Counters   : counters,
		Instances  : instances,
	}, nil
}

// Returns the performance object name or counter name corresponding to the specified index. 
func PdhLookupPerfNameByIndex(MachineName string, idx uint32) (string, error) {
	var (
		machine         uintptr
		NameBuffer      [4096]uint16
		NameBufferSize  = uint32(4096)
	)

	if len(MachineName) != 0 {
		m, _ := syscall.UTF16PtrFromString(MachineName)
		machine = uintptr(unsafe.Pointer(m))
	}

	r, _, _ := procPdhLookupPerfNameByIndex.Call(
		machine, // local machine
		uintptr(idx),
		uintptr(unsafe.Pointer(&NameBuffer[0])),
		uintptr(unsafe.Pointer(&NameBufferSize)))

	if r != 0 {
		return "", errors.New(codeText[r])
	}

	return syscall.UTF16ToString(NameBuffer[:]), nil
}

func PdhLookupPerfIndexByName(MachineName, CounterName string) (uint32, error) {
	var (
		machine uintptr
		idx uint32
	)

	if len(MachineName) != 0 {
		m, _ := syscall.UTF16PtrFromString(MachineName)
		machine = uintptr(unsafe.Pointer(m))
	}

	NameBuffer, _ := syscall.UTF16PtrFromString(CounterName) 

	r, _, _ := procPdhLookupPerfIndexByName.Call(
		machine, // local machine
		uintptr(unsafe.Pointer(NameBuffer)),
		uintptr(unsafe.Pointer(&idx)))

	if r != 0 {
		return 0, errors.New(codeText[r])
	}

	return idx, nil
}

// Removes a counter from a query.
func PdhRemoveCounter(counterHdl syscall.Handle) error {
	r, 	_, _ := procPdhRemoveCounter.Call(uintptr(counterHdl))

	if r != 0 {
		return errors.New(codeText[r])
	}

	return nil
}
