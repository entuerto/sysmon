// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mem

import (
	"fmt"
	"log"
//	"testing"
)


func ExampleSwapMemory() {
	s, err := SwapMemory()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", s)
	// Output: _
}

func ExampleVirtualMemory() {
	v, err := VirtualMemory()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", v)
	// Output: _
}