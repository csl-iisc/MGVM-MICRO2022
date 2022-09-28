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
	// engine.AcceptHook(akita.NewEventLogger(log.New(os.Stdout, "", 0)))
	t := acceptance.NewTest()

	createNetwork(engine, t)
	t.GenerateMsgs(1)
	// t.GenerateMsgs(512)

	engine.Run()

	t.MustHaveReceivedAllMsgs()
	t.ReportBandwidthAchieved(engine.CurrentTime())
	t.ReportAverageLatency()
	atexit.Exit(0)
}

func createNetwork(engine akita.Engine, test *acceptance.Test) {
	freq := 1 * akita.GHz
	var agents []*acceptance.Agent
	numAgents := 2
	for i := 0; i < numAgents; i++ {
		agent := acceptance.NewAgent(
			engine, freq, fmt.Sprintf("Agent%d", i), 1, test)
		agent.TickLater(0)
		agents = append(agents, agent)
	}

	chipConnector := chipnetwork.NewConnector()
	chipConnector = chipConnector.
		WithEngine(engine).
		WithSwitchLatency(144 + 12).
		WithFreq(1 * akita.GHz).
		WithFlitByteSize(64).
		WithNetworkName("ICN").
		WithNumReqPerCycle(12)

	chipConnector.CreateNetwork()

	for i := 0; i < numAgents; i++ {
		chipConnector.PlugInChip(agents[i].Ports)
	}

	// for i := 1; i < 2; i++ {
	// 	chipConnector.PlugInChip(agents[i].Ports)
	// }

	chipConnector.MakeNetwork()

	for i := 0; i < numAgents; i++ {
		test.RegisterAgent(agents[i])
	}
	// test.RegisterAgent(agents[1])
}
