// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"time"

	"github.com/entuerto/sysmon/disk"
)

var (
	header = fmt.Sprintf("  %-20s %10s %10s %10s %10s %10s %10s    \n" ,"name", "rc/s", "wc/s", "rb/s", "wb/s", "Read/t", "Write/t")
	line   = "  %-20s %10d %10d %10s %10s %10s %10s   \n"
)

func delta(f, s *disk.IOCounters) disk.IOCounters {
	return disk.IOCounters{
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

func toDuration(v uint64) time.Duration {
	return time.Duration(v * 100) * time.Nanosecond
}

func main() {
	qio, err := disk.QueryIOCounters("\\\\.\\PhysicalDrive0", time.Second)

	if err != nil {
		log.Println(err)
	}

	fmt.Printf(header)
	go func() {
		pioc := <- qio.IOCounterChan

		for {
			ioc := <- qio.IOCounterChan

			dio := delta(pioc, ioc)
			fmt.Printf(line, 
				       dio.Name, 
				       dio.ReadCount, 
				       dio.WriteCount, 
				       dio.ReadBytes, 
				       dio.WriteBytes, 
				       toDuration(dio.ReadTime), 
				       toDuration(dio.WriteTime))	
			pioc = ioc
		}
	}()

	<-time.After(20 * time.Second)
    qio.Stop()
}
