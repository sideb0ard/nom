package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mgutz/ansi"
)

var colors = []string{ansi.ColorCode("red+b:white"), ansi.ColorCode("white+h:blue")}
var reset = ansi.ColorCode("reset")
var toggle = 1

func main() {
	flag.Parse()
	// ifs, err := pcap.FindAllDevs()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for _, v := range ifs {
	// 	fmt.Println(v.Name)
	// }
	handle, err := pcap.OpenLive("en0", 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for pkt := range packetSource.Packets() {
		//fmt.Println("Pkt is:", reflect.TypeOf(pkt))
		nom(pkt)
		//eth, _ := layers.LayerTypeEthernet.(*layers.Ethernet)
		//fmt.Println(eth.SrcMac, eth.DstMac)
		//if ethLayer := pkt.Layer(layers.LayerTypeEthernet); ethLayer != nil {
		//	fmt.Println("Ethernet")
		//}

		//if ipLayer := pkt.Layer(layers.LayerTypeIPv4); ipLayer != nil {
		//	fmt.Println("IP")
		//	if tcpLayer := pkt.Layer(layers.LayerTypeTCP); tcpLayer != nil {
		//		// Get actual TCP data from this layer
		//		tcp, _ := tcpLayer.(*layers.TCP)
		//		fmt.Printf("From src port %d to dst port %d\n", tcp.SrcPort, tcp.DstPort)
		//	}
		//}
		//nom(pkt)
		//for _, layer := range pkt.Layers() {
		//	//fmt.Println("PACKET LAYER:", layer.LayerType())
		//	fmt.Println("PACKET LAYER:", reflect.TypeOf(layer))
		//}
	}

}

func nom(pkt gopacket.Packet) {
	fmt.Println(colors[toggle], pkt, reset)
	toggle = 1 - toggle
	fmt.Println()
}
