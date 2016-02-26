// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package disk

import (
	"fmt"
	"time"
	"testing"
)

func TestAllPartitions(t *testing.T) {
	partitions, err := AllPartitions()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(partitions) == 0 {
		t.Errorf("could not get partitions")
	}

	for _, p := range partitions {
		fmt.Printf("%#v\n", p)

		if p.DeviceId == "" {
			t.Errorf("could not get Partition information: %v", p)
		}
	}
}

func TestAllVolumes(t *testing.T) {
	volumes, err := AllVolumes()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(volumes) == 0 {
		t.Errorf("could not get volumes")
	}

	for _, v := range volumes {
		fmt.Printf("%#v\n", v)

		if v.DeviceId == "" {
			t.Errorf("could not get Volume information: %v", v)
		}
	}
}

func TestAllUsage(t *testing.T) {
	usages, err := AllUsage()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(usages) == 0 {
		t.Errorf("could not get usage")
	}

	for _, u := range usages {
		fmt.Printf("%#v\n", u)

		if u.DeviceId == "" {
			t.Errorf("could not get usage information: %v", u)
		}
	}
}

func TestAllDrives(t *testing.T) {
	drives, err := AllDrives()

	if err != nil {
		t.Errorf("error %v", err)
	}

	if len(drives) == 0 {
		t.Errorf("could not get drives")
	}

	for _, d := range drives {
		fmt.Printf("%#v\n", d)

		if d.DeviceId == "" {
			t.Errorf("could not get drive information: %v", d)
		}
	}
}

func TestQueryIOCounters(t *testing.T) {
	qio, err := QueryIOCounters("\\\\.\\PhysicalDrive0", time.Second)

	if err != nil {
		t.Errorf("error %v", err)
	}

	go func() {
		pioc := <- qio.IOCounterChan

		for {
			ioc := <- qio.IOCounterChan
			fmt.Printf("%#v\n", delta(pioc, ioc))	
			pioc = ioc
		}
	}()

	<-time.After(20 * time.Second)
    qio.Stop()
}

func delta(f, s *IOCounters) IOCounters {
	return IOCounters{
		Name       :  s.Name,
		ReadCount  :  s.ReadCount - f.ReadCount,
		WriteCount :  s.WriteCount - f.WriteCount,
		ReadBytes  :  s.ReadBytes - f.ReadBytes,
		WriteBytes :  s.WriteBytes - f.WriteBytes,
		ReadTime   :  s.ReadTime - f.ReadTime,
		WriteTime  :  s.WriteTime - f.WriteTime,
		IoTime     :  s.IoTime,
	}
}
