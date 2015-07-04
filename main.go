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

	var stateChannel = make(chan map[string]int)
	go oldTimer(timerChannel)
	go getSocketStatus(stateChannel, timerChannel)

	portdataz := make([]string, len(STATE)+2) // +2 for space and total col
	go updateSocketStatusData(stateChannel, portdataz)

	ls := ui.NewList()
	ls.Border.Label = "Port Status Count"
	ls.Items = portdataz
	ls.Width = 22
	ls.Height = len(portdataz) + 2 // 2 for top/bottom border

	draw := func() {
		ui.Render(ls)
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
