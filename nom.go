package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"

	ui "github.com/gizak/termui"
)

var STATE = map[string]string{
	"01": "ESTABLISHED",
	"02": "SYN_SENT",
	"03": "SYN_RECV",
	"04": "FIN_WAIT1",
	"05": "FIN_WAIT2",
	"06": "TIME_WAIT",
	"07": "CLOSE",
	"08": "CLOSE_WAIT",
	"09": "LAST_ACK",
	"0A": "LISTEN",
	"0B": "CLOSING",
	"0C": "NEW_SYN_RECV",
}

var NETFILES = map[string]string{
	"tcp":  "/proc/net/tcp",
	"tcp6": "/proc/net/tcp6",
}

func reverseMap(m map[string]string) map[string]string {
	n := make(map[string]string)
	for k, v := range m {
		n[v] = k
	}
	return n
}

func getSocketStatus(stateChannel chan map[string]int) map[string]int {
	for {
		portCounts := make(map[string]int)
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

				_, ok := portCounts[state]
				if ok {
					portCounts[state] += 1
				} else {
					portCounts[state] = 1
				}
			}
		}
		stateChannel <- portCounts
		time.Sleep(time.Second / 4)
	}
}

func updateSocketStatusData(stateChannel chan map[string]int, spdataz [][]int) {
	revStates := reverseMap(STATE)
	for {
		select {
		case portCounts := <-stateChannel:
			for state, count := range portCounts {
				nomnum, _ := strconv.ParseInt(revStates[state], 16, 8)
				spdataz[nomnum-1] = append(spdataz[nomnum-1][1:], count)
			}
		default:
		}
	}
}

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.UseTheme("helloworld")

	rand.Seed(time.Now().UnixNano())
	//in := make(chan int)
	//out := make(chan int, 20)
	done := make(chan bool)

	p := ui.NewPar(":PRESS q TO QUIT DEMO")
	p.Height = 3
	p.Width = 50
	p.TextFgColor = ui.ColorWhite
	p.Border.Label = "Socket State Sparklines"
	p.Border.FgColor = ui.ColorCyan

	spdataz := make([][]int, len(STATE))
	for i := 0; i < len(STATE); i++ {
		for j := 0; j < 100; j++ {
			spdataz[i] = append(spdataz[i], 0)
		}
	}

	var stateChannel = make(chan map[string]int)

	go getSocketStatus(stateChannel)
	go updateSocketStatusData(stateChannel, spdataz)

	skeys := make([]string, len(STATE))
	i := 0
	for k, _ := range STATE {
		skeys[i] = k
		i++
	}
	sort.Strings(skeys)

	spStates := []ui.Sparkline{}
	for _, v := range skeys {
		spark := ui.Sparkline{}
		spark.Height = 1
		spark.Title = STATE[v]
		nomnum, _ := strconv.ParseInt(v, 16, 8)
		spark.Data = spdataz[nomnum-1]
		spark.LineColor = ui.ColorCyan
		spark.TitleColor = ui.ColorWhite
		spStates = append(spStates, spark)
	}

	sp := ui.NewSparklines(spStates[0], spStates[1], spStates[2], spStates[3], spStates[4], spStates[5], spStates[6], spStates[7], spStates[8], spStates[9], spStates[10], spStates[11])
	sp.Width = 50
	sp.Height = 25
	sp.Border.Label = "Sparkline"
	sp.Y = 3

	draw := func(t int) {
		for i := 0; i < len(spdataz); i++ {
			sp.Lines[i].Data = spdataz[i]
		}
		ui.Render(p, sp)
	}

	evt := ui.EventCh()
	j := 0
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		default:
			draw(i)
			j++
			time.Sleep(time.Second / 4)
		}
	}

	<-done
}
