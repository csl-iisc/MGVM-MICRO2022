package dispatching

import (
	// "fmt"
	"gitlab.com/akita/mgpusim/kernels"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/mgpusim/timing/cp/internal/resource"
)

// oldChipletAwareAlgorithm can dispatch workgroups to CUs in a round robin
// fasion.
type oldChipletAwareAlgorithm struct {
	gridBuilder kernels.GridBuilder
	cuPool      resource.CUResourcePool
	numChiplets int

	currWG *kernels.WorkGroup
	nextCU int
	// nextChipletID    int
	numDispatchedWGs int
}

// RegisterCU allows the oldChipletAwareAlgorithm to dispatch work-group to the CU.
func (a *oldChipletAwareAlgorithm) RegisterCU(cu resource.DispatchableCU) {
	a.cuPool.RegisterCU(cu)
}

// StartNewKernel lets the algorithms to start dispatching a new kernel.
func (a *oldChipletAwareAlgorithm) StartNewKernel(info kernels.KernelLaunchInfo) {
	a.numDispatchedWGs = 0
	a.gridBuilder.SetKernel(info)
}

// NumWG returns the number of work-groups in the currently-dispatching
// work-group.
func (a *oldChipletAwareAlgorithm) NumWG() int {
	return a.gridBuilder.NumWG()
}

// HasNext check if there are more work-groups to dispatch.
func (a *oldChipletAwareAlgorithm) HasNext() bool {
	return a.numDispatchedWGs < a.gridBuilder.NumWG()
}

// Next finds the location to dispatch the next work-group.
// TODO: WORKS ONLY FOR ONE APPLICATION AS OF NOW!!!
func (a *oldChipletAwareAlgorithm) Next() (location dispatchLocation) {
	if a.currWG == nil {
		a.currWG = a.gridBuilder.NextWG()
	}
	// fmt.Println(a.currWG.IDX, a.currWG.IDY, a.currWG.IDZ)
	a.numChiplets = 4

	cusPerChiplet := a.cuPool.NumCU() / a.numChiplets
	// wgsPerChiplets := a.gridBuilder.NumWG() / a.numChiplets

	// startingChiplet := a.nextChipletID
	// TODO: This is kinda of inefficient.
	for c := 0; c < a.numChiplets; c++ {
		// chiplet := (startingChiplet + c) % a.numChiplets
		for i := 0; i < cusPerChiplet; i++ {
			// cuID := chiplet*cusPerChiplet + ((a.nextCU + i) % cusPerChiplet)
			cuID := c*cusPerChiplet + ((a.nextCU + i) % cusPerChiplet)
			cu := a.cuPool.GetCU(cuID)

			locations, ok := cu.ReserveResourceForWG(a.currWG)
			if ok {
				a.nextCU = (cuID + 1) % cusPerChiplet
				// a.nextChipletID = (chiplet + 1) % a.numChiplets
				dispatch := dispatchLocation{
					valid: true,
					cu:    cu.DispatchingPort(),
					cuID:  cuID,
					wg:    a.currWG,
				}
				dispatch.locations =
					make([]protocol.WfDispatchLocation, len(locations))
				for i, localtion := range locations {
					dispatch.locations[i] = protocol.WfDispatchLocation(localtion)
				}

				a.currWG = nil
				a.numDispatchedWGs++
				return dispatch
			}
		}
	}

	return dispatchLocation{}

}

// FreeResources marks the dispatched location to be available.
func (a *oldChipletAwareAlgorithm) FreeResources(location dispatchLocation) {
	a.cuPool.GetCU(location.cuID).FreeResourcesForWG(location.wg)
}
