package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/intel-go/nff-go/devices"
	"github.com/intel-go/nff-go/flow"
	"github.com/intel-go/nff-go/packet"
)

func main() {
	nic  := "rename4"
	bind := "mlx4_core"
	port := uint16(0)

	config := flow.Config {
	}

	flow.CheckFatal(flow.SystemInit(&config))

	device, err := devices.New(nic)
	flow.CheckFatal(err)

	driver, err := device.CurrentDriver()
	flow.CheckFatal(err)

	defer func() {
		flow.CheckFatal(flow.SystemStop())

		// Re-Bind to original driver
		if driver != bind {
			fmt.Printf("Restoring driver: %s\n", driver)
			device.Bind(driver)
		}
	}()

	// Bind to new user specified driver
	device.Bind(bind)

	mainFlow, err := flow.SetReceiver(port)
	flow.CheckFatal(err)

	flow.CheckFatal(flow.SetHandler(mainFlow, handler, nil))

	flow.CheckFatal(flow.SetStopper(mainFlow))

	go func() {
		flow.CheckFatal(flow.SystemStart())
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	flow.CheckFatal(flow.SystemStop())
}

func handler(packet *packet.Packet, context flow.UserContext) {
	fmt.Printf("Packet: %v", packet)
}
