package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"strings"
)

var networkInfoFile = "/proc/net/dev"

type ifaceTraffic struct {
	BytesInPerSecond  int
	BytesOutPerSecond int
}

func (m *ifaceTraffic) Bytes(bi int, bo int) {
	m.BytesInPerSecond = bi
	m.BytesOutPerSecond = bo
}

func numIface() int {
	infoz, err := ioutil.ReadFile(networkInfoFile)
	if err != nil {
		log.Fatal(err)
	}
	data := strings.Split(string(infoz), "\n")
	data = data[2:] // drop title lines
	return len(data)
}

func getIfaceStatus(ifaceChannel chan map[string]*ifaceTraffic, timerChannel chan int) {
	prevCountIn := make(map[string]int)
	prevCountOut := make(map[string]int)
	ifaceTrafficz := make(map[string]*ifaceTraffic)

	for {
		select {
		case _ = <-timerChannel:
			infoz, err := ioutil.ReadFile(networkInfoFile)
			if err != nil {
				log.Fatal(err)
			}

			data := strings.Split(string(infoz), "\n")
			data = data[2:] // drop title lines

			for _, l := range data {
				if l == "" { // drop empty lines
					continue
				}

				d := strings.Fields(l)
				iface := strings.Split(d[0], ":")[0]
				bytesIn, _ := strconv.Atoi(d[1])
				bytesOut, _ := strconv.Atoi(d[9])

				_, ok := prevCountIn[iface]
				if !ok {
					prevCountIn[iface] = bytesIn
					prevCountOut[iface] = bytesOut
					ifaceTrafficz[iface] = &ifaceTraffic{}
				} else {

					ifaceTrafficz[iface].Bytes((bytesIn-prevCountIn[iface])*2, (bytesOut-prevCountOut[iface])*2) // TODO: fix *2 to not be static (sleep time)
					prevCountIn[iface] = bytesIn
					prevCountOut[iface] = bytesOut
				}
			}
			ifaceChannel <- ifaceTrafficz
		}
	}
}

func updateIfaceStatusData(ifaceChannel chan map[string]*ifaceTraffic, ifacedataz []string) {
	for {
		select {
		case ifaceTrafficz := <-ifaceChannel:
			keysIF := make([]string, 0, len(ifaceTrafficz)*2)
			for key, _ := range ifaceTrafficz {
				keysIF = append(keysIF, key)
			}
			sort.Strings(keysIF)
			i := 0
			for _, key := range keysIF {
				ifacedataz[i] = fmt.Sprintf("%15s : %10d Kb/sec Incoming", key, ifaceTrafficz[key].BytesInPerSecond*8/1024)
				i++
				ifacedataz[i] = fmt.Sprintf("%15s : %10d Kb/sec Outgoing", key, ifaceTrafficz[key].BytesOutPerSecond*8/1024)
				i++
			}

		}
	}
}
