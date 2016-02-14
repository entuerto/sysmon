
// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proc

import (
	"fmt"
	"testing"
)

func TestFindProcess(t *testing.T) {
	p, err := OpenProcess(784)
	if err != nil {
		t.Errorf("error %v", err)
	}

	fmt.Printf("%#v\n", p)

	u, err := p.Usage()
	if err != nil {
		t.Errorf("error %v", err)
	}

	fmt.Printf("%#v\n", u)

	ioc, err := p.IOCounters()
	if err != nil {
		t.Errorf("error %v", err)
	}

	fmt.Printf("%#v\n", ioc)

	pmc, err := p.MemoryInfo()
	if err != nil {
		t.Errorf("error %v", err)
	}

	fmt.Printf("%#v\n", pmc)

	p.Modules()
}
