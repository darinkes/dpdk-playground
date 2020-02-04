package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	portPtr := flag.Int("port", 0, "port number")
	flag.Parse()

	config := flow.Config{
		DPDKArgs: []string {
			"--vdev=net_vdev_netvsc0,iface=eth1",
			"--vdev=net_vdev_netvsc1,iface=eth2",
		},
	}
	flow.CheckFatal(flow.SystemInit(&config))

	port := uint16(*portPtr)
	mainFlow, err := flow.SetReceiver(port)
	flow.CheckFatal(err)

	flow.CheckFatal(flow.SetHandler(mainFlow, handler, nil))
	flow.CheckFatal(flow.SetSender(mainFlow, port))

	go func() {
		flow.CheckFatal(flow.SystemStart())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Stopping...")

	flow.CheckFatal(flow.SystemStop())
}

func handler(packet *packet.Packet, context flow.UserContext) {
	gopacket := gopacket.NewPacket(packet.GetRawPacketBytes(), layers.LayerTypeEthernet, gopacket.Default)
	fmt.Printf("Packet: %v", gopacket)
}
