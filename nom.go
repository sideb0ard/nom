package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/google/gopacket"
	"github.com/mgutz/ansi"
)

var colors = []string{ansi.ColorCode("red+b:white"), ansi.ColorCode("white+h:blue")}
var reset = ansi.ColorCode("reset")
var toggle = 1

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

type tcpSocketStatus struct {
	index         int
	local_address string
	rem_address   string
	status        string
	tx_q_rx_q     string
	timer_when    string
	retransmit    string
	uid           int
	timeout       int
	inode         int
}

var NETFILES = map[string]string{
	"tcp":  "/proc/net/tcp",
	"tcp6": "/proc/net/tcp6",
}

// func getSocketStatus() []tcpSocketStatus {
func getSocketStatus() {

	portCounts := make(map[string]int)

	for proto := range NETFILES {
		infoz, err := ioutil.ReadFile(NETFILES[proto])
		if err != nil {
			log.Fatal(err)
		}

		data := strings.Split(string(infoz), "\n")
		data = data[1:]
		for _, l := range data {
			if l == "" {
				continue
			}
			d := strings.Fields(l)
			//idx, _ := strconv.Atoi(d[0])
			state := STATE[d[3]]

			_, ok := portCounts[state]
			if ok {
				portCounts[state] += 1
			} else {
				portCounts[state] = 1
			}

			//ts := tcpSocketStatus{index: idx, local_address: d[1], rem_address: d[2], status: STATE[d[3]]}
			//fmt.Println(ts)

		}
	}
	fmt.Println(portCounts)
}

func main() {
	getSocketStatus()
	// flag.Parse()
	// ifs, err := pcap.FindAllDevs()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, v := range ifs {
	// 	fmt.Println(v.Name)
	// }
	//handle, err := pcap.OpenLive("en0", 65536, true, pcap.BlockForever)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer handle.Close()

	//packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	//for pkt := range packetSource.Packets() {
	//	//fmt.Println("Pkt is:", reflect.TypeOf(pkt))
	//	nom(pkt)
	//	//eth, _ := layers.LayerTypeEthernet.(*layers.Ethernet)
	//	//fmt.Println(eth.SrcMac, eth.DstMac)
	//	//if ethLayer := pkt.Layer(layers.LayerTypeEthernet); ethLayer != nil {
	//	//	fmt.Println("Ethernet")
	//	//}

	//	//if ipLayer := pkt.Layer(layers.LayerTypeIPv4); ipLayer != nil {
	//	//	fmt.Println("IP")
	//	//	if tcpLayer := pkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
	//	//		// Get actual TCP data from this layer
	//	//		tcp, _ := tcpLayer.(*layers.TCP)
	//	//		fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
	//	//	}
	//	//}
	//	//nom(pkt)
	//	//for _, layer := range pkt.Layers() {
	//	//	//fmt.Println("PACKET LAYER:", layer.LayerType())
	//	//	fmt.Println("PACKET LAYER:", reflect.TypeOf(layer))
	//	//}
	//}

}

func nom(pkt gopacket.Packet) {
	fmt.Println(colors[toggle], pkt, reset)
	toggle = 1 - toggle
	fmt.Println()
}
