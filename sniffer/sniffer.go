package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/intel-go/nff-go/devices"
	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func main() {
	nic  := "0002:00:02.0"
	bind := "mlx4_core"
	port := uint16(0)

	flow.CheckFatal(flow.SystemInit(nil))

	device, err := devices.New(nic)
	flow.CheckFatal(err)

	driver, err := device.CurrentDriver()
	flow.CheckFatal(err)

	device.Bind(bind)

	loadedDriver, _ := device.CurrentDriver()
	fmt.Printf("Driver: %v\n", loadedDriver)

	mainFlow, err := flow.SetReceiver(port)
	flow.CheckFatal(err)

	flow.CheckFatal(flow.SetHandler(mainFlow, handler, nil))
	flow.CheckFatal(flow.SetSender(mainFlow, 0))

	go func() {
		flow.CheckFatal(flow.SystemStart())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	fmt.Println("Stopping...")

	flow.CheckFatal(flow.SystemStop())

	if driver != bind {
		fmt.Printf("Restoring driver: %s\n", driver)
		device.Bind(driver)
	}
}

func handler(packet *packet.Packet, context flow.UserContext) {
	gopacket := gopacket.NewPacket(packet.GetRawPacketBytes(), layers.LayerTypeEthernet, gopacket.Default)
	fmt.Printf("Packet: %v", gopacket)
}
