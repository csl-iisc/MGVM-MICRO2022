package tracing_test

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/util/tracing"
)

type SampleDomain struct {
	*akita.HookableBase
	taskIDs []int
	nextID  int
}

func (d *SampleDomain) Name() string {
	return "sample domain"
}

func (d *SampleDomain) Start(now akita.VTimeInSec) {
	tracing.StartTask(
		fmt.Sprintf("%d", d.nextID),
		"",
		now,
		d,
		"sampleTaskKind",
		"something",
		nil,
	)
	d.taskIDs = append(d.taskIDs, d.nextID)

	d.nextID++
}

func (d *SampleDomain) End(now akita.VTimeInSec) {
	tracing.EndTask(
		fmt.Sprintf("%d", d.taskIDs[0]),
		now,
		d,
	)
	d.taskIDs = d.taskIDs[1:]
}

// Example for how to use standard tracers
func ExampleTracer() {
	domain := &SampleDomain{
		HookableBase: akita.NewHookableBase(),
	}

	filter := func(t tracing.Task) bool {
		return t.Kind == "sampleTaskKind"
	}

	totalTimeTracer := tracing.NewTotalTimeTracer(filter)
	busyTimeTracer := tracing.NewBusyTimeTracer(filter)
	avgTimeTracer := tracing.NewAverageTimeTracer(filter)
	tracing.CollectTrace(domain, totalTimeTracer)
	tracing.CollectTrace(domain, busyTimeTracer)
	tracing.CollectTrace(domain, avgTimeTracer)

	domain.Start(1)
	domain.Start(1.5)
	domain.End(2)
	domain.End(3)

	fmt.Println(totalTimeTracer.TotalTime())
	fmt.Println(busyTimeTracer.BusyTime())
	fmt.Println(avgTimeTracer.AverageTime())

	// Output:
	// 2.5
	// 2
	// 1.25
}
