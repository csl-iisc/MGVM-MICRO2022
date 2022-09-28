package dispatching

import (
	"fmt"

	"gitlab.com/akita/mgpusim/kernels"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/mgpusim/timing/cp/internal/resource"
)

// chipletAwareAlgorithm can dispatch workgroups to CUs in a round robin
// fasion.
type chipletAwareAlgorithm struct {
	gridBuilder kernels.GridBuilder
	cuPool      resource.CUResourcePool
	numChiplets int

	currWG           *kernels.WorkGroup
	nextCU           []int
	nextChipletID    int
	numDispatchedWGs int

	perChipletListOfWG [][]*kernels.WorkGroup
}

// RegisterCU allows the chipletAwareAlgorithm to dispatch work-group to the CU.
func (a *chipletAwareAlgorithm) RegisterCU(cu resource.DispatchableCU) {
	a.cuPool.RegisterCU(cu)
}

// StartNewKernel lets the algorithms to start dispatching a new kernel.
func (a *chipletAwareAlgorithm) StartNewKernel(info kernels.KernelLaunchInfo) {
	a.numDispatchedWGs = 0
	a.gridBuilder.SetKernel(info)

	a.numChiplets = 4

	a.nextCU = make([]int, a.numChiplets)

	a.perChipletListOfWG = make([][]*kernels.WorkGroup, a.numChiplets)
	for i := 0; i < a.numChiplets; i++ {
		a.perChipletListOfWG[i] = make([]*kernels.WorkGroup, 0)
		a.nextCU[i] = 0
	}

	// This is a do-while loop
	a.currWG = a.gridBuilder.NextWG()
	divisor := 0
	wgCoordinate := ""
	if a.currWG.Packet.GridSizeZ > 1 {
		divisor = int(a.currWG.Packet.GridSizeZ) / int(a.currWG.Packet.WorkgroupSizeZ*uint16(a.numChiplets))
		wgCoordinate = "Z"
	} else if a.currWG.Packet.GridSizeY > 1 {
		divisor = int(a.currWG.Packet.GridSizeY) / int(a.currWG.Packet.WorkgroupSizeY*uint16(a.numChiplets))
		wgCoordinate = "Y"
	} else {
		divisor = int(a.currWG.Packet.GridSizeX) / int(a.currWG.Packet.WorkgroupSizeX*uint16(a.numChiplets))
		wgCoordinate = "X"
	}
	for a.currWG != nil {
		fmt.Println(a.currWG.Packet.GridSizeX, a.currWG.Packet.WorkgroupSizeX, a.numChiplets)
		wgID := 0
		switch wgCoordinate {
		case "X":
			wgID = a.currWG.IDX
		case "Y":
			wgID = a.currWG.IDY
		case "Z":
			wgID = a.currWG.IDZ
		}
		x := 0
		if divisor == 0 {
			x = 0
		} else {
			x = wgID / divisor
		}
		if x >= a.numChiplets {
			x = 0
		}

		a.perChipletListOfWG[x] = append(a.perChipletListOfWG[x], a.currWG)
		a.currWG = a.gridBuilder.NextWG()
	}

}

// NumWG returns the number of work-groups in the currently-dispatching
// work-group.
func (a *chipletAwareAlgorithm) NumWG() int {
	return a.gridBuilder.NumWG()
}

// HasNext check if there are more work-groups to dispatch.
func (a *chipletAwareAlgorithm) HasNext() bool {
	return a.numDispatchedWGs < a.gridBuilder.NumWG()
}

// Next finds the location to dispatch the next work-group.
// TODO: WORKS ONLY FOR ONE APPLICATION AS OF NOW!!!
func (a *chipletAwareAlgorithm) Next() (location dispatchLocation) {
	a.numChiplets = 4

	cusPerChiplet := a.cuPool.NumCU() / a.numChiplets
	// wgsPerChiplets := a.gridBuilder.NumWG() / a.numChiplets

	startingChiplet := a.nextChipletID
	// TODO: This is kinda of inefficient.
	for c := 0; c < a.numChiplets; c++ {
		chiplet := (startingChiplet + c) % a.numChiplets

		// a.currWG = nil

		if len(a.perChipletListOfWG[chiplet]) != 0 {
			a.currWG = a.perChipletListOfWG[chiplet][0]
		} else {
			continue
		}
		// for len(a.perChipletListOfWG[c]) == 0 {
		// 	a.currWG = a.gridBuilder.NextWG()
		// 	if a.currWG == nil {
		// 		break
		// 	}
		// 	x := a.currWG.IDX / (a.currWG.SizeX / a.numChiplets)
		// 	a.perChipletListOfWG[x] = append(a.perChipletListOfWG[x], a.currWG)
		// }

		if a.currWG == nil {
			continue
		}

		for i := 0; i < cusPerChiplet; i++ {
			// cuID := chiplet*cusPerChiplet + ((a.nextCU + i) % cusPerChiplet)
			cuID := chiplet*cusPerChiplet + ((a.nextCU[chiplet] + i) % cusPerChiplet)
			cu := a.cuPool.GetCU(cuID)

			locations, ok := cu.ReserveResourceForWG(a.currWG)
			if ok {
				a.nextCU[chiplet] = (cuID + 1) % cusPerChiplet
				a.nextChipletID = (chiplet + 1) % a.numChiplets
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

				fmt.Println(a.currWG.IDX, a.currWG.IDY, a.currWG.IDZ, chiplet)

				a.currWG = nil
				a.perChipletListOfWG[chiplet] = a.perChipletListOfWG[chiplet][1:]
				a.numDispatchedWGs++

				return dispatch
			} else {
				// cu.FreeResourcesForWG(a.currWG)
				// a.perChipletListOfWG[c] = append(a.perChipletListOfWG[c], a.currWG)
				// a.currWG = nil
			}
		}
	}

	return dispatchLocation{}

}

// FreeResources marks the dispatched location to be available.
func (a *chipletAwareAlgorithm) FreeResources(location dispatchLocation) {
	a.cuPool.GetCU(location.cuID).FreeResourcesForWG(location.wg)
}
