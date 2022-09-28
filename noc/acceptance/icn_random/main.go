package main

import (
	"flag"
	"fmt"
	"math/rand"

	"github.com/tebeka/atexit"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/noc/acceptance"
	"gitlab.com/akita/noc/networking/chipnetwork"
)

func main() {
	flag.Parse()
	rand.Seed(1)

	engine := akita.NewSerialEngine()
	// engine := akita.NewParallelEngine()
	t := acceptance.NewTest()

	createNetwork(engine, t)
	t.GenerateMsgs(2048)

	engine.Run()

	t.MustHaveReceivedAllMsgs()
	t.ReportBandwidthAchieved(engine.CurrentTime())
	atexit.Exit(0)
}

func createNetwork(engine akita.Engine, test *acceptance.Test) {
	freq := 1.0 * akita.GHz
	var agents []*acceptance.Agent
	for i := 0; i < 8; i++ {
		agent := acceptance.NewAgent(
			engine, freq, fmt.Sprintf("Agent%d", i), 1, test)
		agent.TickLater(0)
		agents = append(agents, agent)
		test.RegisterAgent(agent)
	}

	chipConnector := chipnetwork.NewConnector()
	chipConnector = chipConnector.
		WithEngine(engine).
		WithSwitchLatency(32).
		WithFreq(1 * akita.GHz).
		WithNetworkName("ICN")

	chipConnector.CreateNetwork()

	for i := 0; i < 8; i++ {
		chipConnector.PlugInChip(agents[i].Ports)
	}

	chipConnector.MakeNetwork()

}
