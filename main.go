package main

import (
	"fmt"

	"github.com/intel-go/nff-go/devices"
	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"
)

func main() {
	nic  := "enP2s2"
	bind := "mlx4_core"
	port := uint16(0)

	flow.CheckFatal(flow.SystemInit(nil))

	device, err := devices.New(nic)
	flow.CheckFatal(err)

	driver, err := device.CurrentDriver()
	flow.CheckFatal(err)

	defer func() {
		fmt.Printf("Restoring driver: %s\n", driver)
		flow.SystemStop()

		// Re-Bind to original driver
		device.Bind(driver)
	}()

	// Bind to new user specified driver
	device.Bind(bind)

	mainFlow, err := flow.SetReceiver(port)
	flow.CheckFatal(err)

	err = flow.SetHandlerDrop(mainFlow, handler, nil)
	flow.CheckFatal(err)

	err = flow.SetSender(mainFlow, port)
	flow.CheckFatal(err)

	err = flow.SystemStart()
	flow.CheckFatal(err)
}

func handler(packet *packet.Packet, context flow.UserContext) bool {
	fmt.Printf("Packet: %v", packet)
	return true
}
