// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package win32

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func ExamplePdhUCpuUsageCounter() {
	h, err := PdhOpenQuery() 
	if err != nil {
		log.Fatal("PdhOpenQuery: ", err)
	}

	counter, err := PdhAddCounter(h, "\\Processeur(_Total)\\% temps processeur")
	if err != nil {
		log.Fatal("PdhAddCounter: ", err)
	}

	if err := PdhCollectQueryData(h); err != nil {
		log.Fatal("PdhCollectQueryData: ", err)
	}

	go func() {
		for {
			time.Sleep(time.Second)

			if err := PdhCollectQueryData(h); err != nil {
				log.Fatal("PdhCollectQueryData: ", err)
			}

			v, err := PdhGetCounterValueFloat64(counter)
			if err != nil {
				log.Fatal("PdhGetCounterValueFloat64: ", err)
			} 
	
			fmt.Printf("\\Processeur(_Total)\\%% temps processeur : %.2f\n", v)
		}
	}()

	<-time.After(10 * time.Second)

	PdhCloseQuery(h) 
	if err != nil {
		log.Fatal("PdhCloseQuery: ", err)
	}  
	// Output: _
}

func ExamplePdhEnumObjects() {
	objects, err := PdhEnumObjects("", "") 
	if err != nil {
		log.Fatal("PdhEnumObjects: ", err)
	}

	for _, o := range objects {
		fmt.Println(o)
	}
	// Output: _
}

func ExamplePdhEnumObjectItems() {
	obItems, err := PdhEnumObjectItems("", "", "Processeur") 
	if err != nil {
		log.Fatal("PdhEnumObjects: ", err)
	}

	fmt.Println("Counters for ", obItems.ObjectName)
	for _, c := range obItems.Counters {
		fmt.Println(c)
	}
	// Output: _
}
/*
func TestRegQuery(t *testing.T) {
	var (
		valType uint32
		buffer  [231709]uint16
		bufLen  uint32
	)

	bufLen = 231708

	name, _ := syscall.UTF16PtrFromString("Counter 009") 

	err := syscall.RegQueryValueEx(
		syscall.HKEY_PERFORMANCE_DATA, 
		name, 
		nil, 
		&valType, 
		(*byte)(unsafe.Pointer(&buffer[0])), 
		&bufLen)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("buffer: %s \n\n", syscall.UTF16ToString(buffer[:]))
	fmt.Printf("buffer Len: %d \n\n", bufLen)

	start := 0
	for i := 0; i < int(bufLen / 2); i++ {
		if buffer[i] == 0 {
			idx := syscall.UTF16ToString(buffer[start:i])
			start = i + 1

			var desc string
			for i = start; i < int(bufLen / 2); i++ {
				if buffer[i] == 0 {
					desc = syscall.UTF16ToString(buffer[start:i])
					start = i + 1
					break
				}
			}
	        fmt.Printf("Index : %s   -   Text: %s\n", idx, desc)
		}
	}
}
*/
func TestByIndex(t *testing.T) {
	name, err := PdhLookupPerfNameByIndex("", 6) 
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Found Index : %d   -   Text: %s\n", 6, name)
}

func TestByName(t *testing.T) {
	idx, err := PdhLookupPerfIndexByName("", "% temps processeur") 
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("Found Index : %d   -   Text: %s\n", idx, "% temps processeur")
}
