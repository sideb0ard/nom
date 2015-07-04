package main

import ui "github.com/gizak/termui"

func main() {

	err := ui.Init()
	if err != nil {
		panic(err)
	}
	defer ui.Close()
	ui.UseTheme("helloworld")

	var timerChannel = make(chan int)
	go oldTimer(timerChannel)

	// Socket Status - two go routines - one to get stats,
	// other to update data structure for use in widgets
	var stateChannel = make(chan map[string]int)
	go getSocketStatus(stateChannel, timerChannel)
	portdataz := make([]string, len(STATE)+2) // +2 for space and total col
	go updateSocketStatusData(stateChannel, portdataz)

	ls := ui.NewList()
	ls.Border.Label = " Port Status Count"
	ls.Items = portdataz
	// ls.Width = 22
	ls.Height = len(portdataz) + 2 // 2 for top/bottom border

	// Inteface Stats - two go routines again, same deal
	var ethyChannel = make(chan []map[string]int)
	go getIfaceStatus(ethyChannel, timerChannel)
	// ifaceDataz := make([]map[string]int)
	ifaceDataz := make([]string, numIface()*2) // 2 for in and out
	go updateIfaceStatusData(ethyChannel, ifaceDataz)

	ethy := ui.NewList()
	ethy.Border.Label = " Interface Traffic Levels"
	ethy.Items = ifaceDataz
	ethy.Width = 70
	ethy.Height = 10

	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(3, 0, ls),
			ui.NewCol(3, 0, ethy)))

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
