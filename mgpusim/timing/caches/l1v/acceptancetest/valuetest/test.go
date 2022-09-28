package main

import (
	"sync"

	"gitlab.com/akita/mem/idealmemcontroller"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mgpusim/timing/caches/l1v"
	"gitlab.com/akita/mem"
	"gitlab.com/akita/mem/acceptancetests"
	"gitlab.com/akita/mem/cache"
)

type test struct {
	engine          akita.Engine
	conn            *akita.DirectConnection
	agent           *acceptancetests.MemAccessAgent
	lowModuleFinder *cache.SingleLowModuleFinder
	dram            *idealmemcontroller.Comp
	c               *l1v.Cache
}

func (t *test) run(wg *sync.WaitGroup) {
	defer wg.Done()

	t.agent.TickLater(0)
	t.engine.Run()
}

func (t *test) setMaxAddr(addr uint64) {
	t.agent.MaxAddress = addr
}

func newTest(name string) *test {
	t := new(test)

	t.engine = akita.NewSerialEngine()
	t.conn = akita.NewDirectConnection("conn", t.engine, 1*akita.GHz)

	t.dram = idealmemcontroller.New("dram", t.engine, 1*mem.GB)
	t.lowModuleFinder = new(cache.SingleLowModuleFinder)
	t.lowModuleFinder.LowModule = t.dram.ToTop

	t.c = l1v.NewBuilder().
		WithEngine(t.engine).
		WithLowModuleFinder(t.lowModuleFinder).
		Build("cache")

	t.agent = acceptancetests.NewMemAccessAgent(t.engine)
	t.agent.WriteLeft = 1000
	t.agent.ReadLeft = 1000
	t.agent.LowModule = t.c.TopPort

	t.conn.PlugIn(t.agent.ToMem, 4)
	t.conn.PlugIn(t.c.TopPort, 4)
	t.conn.PlugIn(t.c.BottomPort, 16)
	t.conn.PlugIn(t.c.ControlPort, 1)
	t.conn.PlugIn(t.dram.ToTop, 16)

	return t
}

func main() {
	var wg sync.WaitGroup

	//rand.Seed(1)

	t1 := newTest("Max_64")
	t1.setMaxAddr(64)
	wg.Add(1)

	t2 := newTest("Max_1024")
	t2.setMaxAddr(1024)
	wg.Add(1)

	t3 := newTest("Max_1M")
	t3.setMaxAddr(1048576)
	wg.Add(1)

	go t1.run(&wg)
	go t2.run(&wg)
	go t3.run(&wg)
	wg.Wait()
}
