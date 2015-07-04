package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var networkInfoFile = "/proc/net/dev"

func getIfaceStatus(ifaceStatusChannel chan []map[string]int, timerChannel chan int) {
	prevCountIn := make(map[string]int)
	prevCountOut := make(map[string]int)
	curCountIn := make(map[string]int)
	curCountOut := make(map[string]int)
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
				} else {
					curCountIn[iface] = (bytesIn - prevCountIn[iface]) * 2    // sleep time
					curCountOut[iface] = (bytesOut - prevCountOut[iface]) * 2 // sleep time
					prevCountIn[iface] = bytesIn
					prevCountOut[iface] = bytesOut
				}

				fmt.Println(iface, ":", curCountIn[iface], "bytes/s", curCountOut[iface], "bytes/s")
			}
			ifaceStatusChannel <- []map[string]int{curCountIn, curCountOut}
		}
	}
}

// func updateSocketStatusData(stateChannel chan map[string]int, dataz []string) {
// 	revStates := reverseMap(STATE)
// 	totalcount := 0
// 	for {
// 		select {
// 		case portCounts := <-stateChannel:
// 			totalcount = 0
// 			for state, count := range portCounts {
// 				nomnum, _ := strconv.ParseInt(revStates[state], 16, 8)
// 				dataz[nomnum-1] = fmt.Sprintf(state+": %d", count)
// 				totalcount += count
// 			}
// 			dataz[len(dataz)-2] = ""
// 			dataz[len(dataz)-1] = fmt.Sprintf("TOTAL       : %d", totalcount)
// 		}
// 	}
// }
