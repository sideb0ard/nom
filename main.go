package main

import (
	"time"

	ui "github.com/gizak/termui"
)

const (
	SLEEPTIME = time.Second / 2
)

func main() {

	// SETUP
	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.UseTheme("helloworld")

	var timerChannel = make(chan int)
	go oldTimer(timerChannel)

	// DATA WORKERS ////////////////////////////

	// Socket Status - two go routines - one to get stats,
	// other to update data structure for use in widgets
	var socketChannel = make(chan map[string]int)
	go getSocketStatus(socketChannel, timerChannel)
	socketdataz := make([]string, len(STATE)+2) // +2 for space and total col
	go updateSocketStatusData(socketChannel, socketdataz)

	// Interface Traffic Stats
	var ifaceChannel = make(chan map[string]*ifaceTraffic)
	go getIfaceStatus(ifaceChannel, timerChannel)
	ifacedataz := make([]string, numIface()*2)
	go updateIfaceStatusData(ifaceChannel, ifacedataz)

	// UI Widgets ////////////////////////////

	sockets := ui.NewList()
	sockets.Border.Label = " Port Status Count"
	sockets.Items = socketdataz
	// ls.Width = 22
	sockets.Height = len(socketdataz) + 2 // 2 for top/bottom border

	ifaceTraffic := ui.NewList()
	ifaceTraffic.Border.Label = " Interface Traffic"
	ifaceTraffic.Items = ifacedataz
	// ls.Width = 22
	ifaceTraffic.Height = len(ifacedataz)

	// BODY Display
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(3, 0, sockets),
			ui.NewCol(3, 0, ifaceTraffic)),
		ui.NewRow(
			ui.NewCol(3, 0, ifaceTraffic)))

	draw := func() {
		ui.Body.Align()
		ui.Render(ui.Body)
	}

	evt := ui.EventCh()
	for {
		select {
		case e := <-evt:
			if e.Type == ui.EventKey && e.Ch == 'q' {
				return
			}
		case _ = <-timerChannel:
			draw()
		}
	}
}
