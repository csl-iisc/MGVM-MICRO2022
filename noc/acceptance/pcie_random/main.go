package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/tebeka/atexit"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/acceptance"
	"gitlab.com/akita/noc/networking/pcie"
)

func main() {
	flag.Parse()
	rand.Seed(1)

	engine := akita.NewSerialEngine()
	t := acceptance.NewTest()

	createNetwork(engine, t)
	t.GenerateMsgs(1000)

	engine.Run()

	t.MustHaveReceivedAllMsgs()
	t.ReportBandwidthAchieved(engine.CurrentTime())
	atexit.Exit(0)
}

func createNetwork(engine akita.Engine, test *acceptance.Test) {
	freq := 1.0 * akita.GHz
	var agents []*acceptance.Agent
	for i := 0; i < 9; i++ {
		agent := acceptance.NewAgent(
			engine, freq, fmt.Sprintf("Agent%d", i), 5, test)
		agent.TickLater(0)
		agents = append(agents, agent)
		test.RegisterAgent(agent)
	}

	pcieConnector := pcie.NewConnector()
	pcieConnector = pcieConnector.
		WithEngine(engine).
		WithNetworkName("PCIe").
		WithVersion3().
		WithX16()

	pcieConnector.CreateNetwork()
	rootComplexID := pcieConnector.CreateRootComplex(agents[0].Ports)
	switch1ID := pcieConnector.AddSwitch(rootComplexID)
	for i := 1; i < 5; i++ {
		pcieConnector.PlugInDevice(switch1ID, agents[i].Ports)
	}

	switch2ID := pcieConnector.AddSwitch(rootComplexID)
	for i := 5; i < 9; i++ {
		pcieConnector.PlugInDevice(switch2ID, agents[i].Ports)
	}
}
