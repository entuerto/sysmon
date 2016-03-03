// Copyright 2015 The sysmon Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package net

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestInterfaces(t *testing.T) {
	ifaces, err := net.Interfaces()

	if err != nil {
		t.Errorf("error %v", err)
	}

	for _, i := range ifaces {
		fmt.Println("Interface{")
		fmt.Println(" Index        : ", i.Index)
		fmt.Println(" Name         : ", i.Name)
		fmt.Println(" MTU          : ", i.MTU)
		fmt.Println(" HardwareAddr : ", i.HardwareAddr)
		fmt.Println(" Flags        : ", i.Flags)
		fmt.Println(" MulticastAddrs: [")

		addrs, _ := i.Addrs()
		for _, a := range addrs {
			fmt.Println("     Addr{")
			fmt.Println("       Network : ", a.Network())
			fmt.Println("       String  : ", a.String())
			fmt.Println("     }")
		}

		multiAddrs, _ := i.MulticastAddrs()
		for _, a := range multiAddrs {
			fmt.Println("     Multicast Addr{")
			fmt.Println("       Network : ", a.Network())
			fmt.Println("       String  : ", a.String())
			fmt.Println("     }")
		}

		fmt.Println("   ]")
		fmt.Println("}")
	}
}

func TestInterfaceAddrs(t *testing.T) {
	addrs, err := net.InterfaceAddrs() 

	if err != nil {
		t.Errorf("error %v", err)
	}

	for _, a := range addrs {
		fmt.Println("Addr{")
		fmt.Println(" Network : ", a.Network())
		fmt.Println(" String  : ", a.String())
		fmt.Println("}")
	}
}

func TestQueryConnections(t *testing.T) {
	conns, err := QueryConnectionsAll() 

	if err != nil {
		t.Errorf("error %v", err)
	}

	for _, c := range conns {
		fmt.Println("Connection{")
		fmt.Printf("  State  : %+v\n", c.State)
		fmt.Printf("  Type   : %+v\n", c.Type)
		fmt.Printf("  Local  : %+v\n", c.Local)
		fmt.Printf("  Remote : %+v\n", c.Remote)
		fmt.Printf("  Pid    : %+v\n", c.Pid)
		fmt.Println("}")
	}
}

func TestQueryIOCounters(t *testing.T) {
	ifaces, err := net.Interfaces()

	if err != nil {
		t.Errorf("error %v", err)
	}

	qio, err := QueryIOCounters(ifaces[0], time.Second)

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
		Name        : s.Name,
		BytesSent   : s.BytesSent - f.BytesSent,
		BytesRecv   : s.BytesRecv - f.BytesRecv,
		PacketsSent : s.PacketsSent - f.PacketsSent,
		PacketsRecv : s.PacketsRecv - f.PacketsRecv,
		Errin       : s.Errin - f.Errin,
		Errout      : s.Errout - f.Errout,
		Dropin      : s.Dropin - f.Dropin,
		Dropout     : s.Dropout - f.Dropin,
	}
}
