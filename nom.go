package main

import (
	"flag"
	"log"

	"code.google.com/p/gopacket/dumpcommand"
	"code.google.com/p/gopacket/pcap"
)

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

	dumpcommand.Run(handle)
}
