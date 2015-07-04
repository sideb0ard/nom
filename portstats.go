package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var STATE = map[string]string{
	"01": "ESTABLISHED ",
	"02": "SYN_SENT    ",
	"03": "SYN_RECV    ",
	"04": "FIN_WAIT1   ",
	"05": "FIN_WAIT2   ",
	"06": "TIME_WAIT   ",
	"07": "CLOSE       ",
	"08": "CLOSE_WAIT  ",
	"09": "LAST_ACK    ",
	"0A": "LISTEN      ",
	"0B": "CLOSING     ",
	"0C": "NEW_SYN_RECV",
}

var NETFILES = map[string]string{
	"tcp":  "/proc/net/tcp",
	"tcp6": "/proc/net/tcp6",
}

func getSocketStatus(stateChannel chan map[string]int, timerChannel chan int) {
	portCounts := make(map[string]int)
	for {
		select {
		case _ = <-timerChannel:
			for _, s := range STATE {
				portCounts[s] = 0
			}
			for proto := range NETFILES {
				infoz, err := ioutil.ReadFile(NETFILES[proto])
				if err != nil {
					log.Fatal(err)
				}

				data := strings.Split(string(infoz), "\n")
				data = data[1:] // drop title line

				for _, l := range data {
					if l == "" { // drop empty lines
						continue
					}

					d := strings.Fields(l)
					state := STATE[d[3]]
					portCounts[state] += 1

				}
			}
			stateChannel <- portCounts
		}
	}
}

func updateSocketStatusData(stateChannel chan map[string]int, dataz []string) {
	revStates := reverseMap(STATE)
	totalcount := 0
	for {
		select {
		case portCounts := <-stateChannel:
			totalcount = 0
			for state, count := range portCounts {
				nomnum, _ := strconv.ParseInt(revStates[state], 16, 8)
				dataz[nomnum-1] = fmt.Sprintf(state+": %d", count)
				totalcount += count
			}
			dataz[len(dataz)-2] = ""
			dataz[len(dataz)-1] = fmt.Sprintf("TOTAL       : %d", totalcount)
		}
	}
}
