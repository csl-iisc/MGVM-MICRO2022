package dispatching

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"gitlab.com/akita/mgpusim/kernels"
	"gitlab.com/akita/mgpusim/protocol"
	"gitlab.com/akita/mgpusim/timing/cp/internal/resource"
)

type chipletLASPAlgorithm struct {
	gridBuilder kernels.GridBuilder
	cuPool      resource.CUResourcePool
	numChiplets int

	currWG           *kernels.WorkGroup
	nextCU           []int
	nextChipletID    int
	numDispatchedWGs int

	perChipletListOfWG [][]*kernels.WorkGroup

	partition string
}

func (a *chipletLASPAlgorithm) RegisterCU(cu resource.DispatchableCU) {
	a.cuPool.RegisterCU(cu)
}

func (a *chipletLASPAlgorithm) StartNewKernel(info kernels.KernelLaunchInfo) {
	a.numDispatchedWGs = 0
	a.gridBuilder.SetKernel(info)

	a.numChiplets = 4

	a.nextCU = make([]int, a.numChiplets)

	a.perChipletListOfWG = make([][]*kernels.WorkGroup, a.numChiplets)
	for i := 0; i < a.numChiplets; i++ {
		a.perChipletListOfWG[i] = make([]*kernels.WorkGroup, 0)
		a.nextCU[i] = 0
	}

	divisor := 0
	mod := 0
	wgCoordinate := ""

	if strings.Contains(a.partition, "Xdiv") {
		divisor = int(math.Ceil(float64(a.gridBuilder.NumWGX()) / float64(a.numChiplets)))
		wgCoordinate = "Xdiv"
	}
	if strings.Contains(a.partition, "Ydiv") {
		divisor = a.gridBuilder.NumWGY() / a.numChiplets
		wgCoordinate = "Ydiv"
	}
	if strings.Contains(a.partition, "Xblk") {
		r, _ := regexp.Compile("Xblk.[0-9]*")
		matched := r.FindString(a.partition)
		v, _ := regexp.Compile("[0-9]+")
		value := v.FindString(matched)
		mod, _ = strconv.Atoi(value)
		wgCoordinate = "Xblk"
	}
	if strings.Contains(a.partition, "Yblk") {
		r, _ := regexp.Compile("Yblk.[0-9]*")
		matched := r.FindString(a.partition)
		v, _ := regexp.Compile("[0-9]+")
		value := v.FindString(matched)
		mod, _ = strconv.Atoi(value)
		wgCoordinate = "Yblk"
	}
	fmt.Println(divisor, mod, wgCoordinate)

	// This is a do-while loop
	a.currWG = a.gridBuilder.NextWG()
	// fmt.Println(a.currWG.Packet.GridSizeX, a.currWG.Packet.WorkgroupSizeX, a.currWG.Packet.GridSizeY, a.currWG.Packet.WorkgroupSizeY, a.numChiplets)

	for a.currWG != nil {
		wgID := 0
		chiplet := 0
		switch wgCoordinate {
		case "Xdiv":
			wgID = a.currWG.IDX
			chiplet = wgID / divisor
		case "Ydiv":
			wgID = a.currWG.IDY
			chiplet = wgID / divisor
		case "Xblk":
			wgID = a.currWG.IDX
			chiplet = (wgID / mod) % a.numChiplets
		case "Yblk":
			wgID = a.currWG.IDY
			chiplet = (wgID / mod) % a.numChiplets
		default:
			panic("paniiiiiccc!!!")
		}
		a.perChipletListOfWG[chiplet] = append(a.perChipletListOfWG[chiplet], a.currWG)
		fmt.Println(a.currWG.IDX, a.currWG.IDY, chiplet)
		a.currWG = a.gridBuilder.NextWG()
	}
}

// NumWG returns the number of work-groups in the currently-dispatching
// work-group.
func (a *chipletLASPAlgorithm) NumWG() int {
	return a.gridBuilder.NumWG()
}

// HasNext check if there are more work-groups to dispatch.
func (a *chipletLASPAlgorithm) HasNext() bool {
	return a.numDispatchedWGs < a.gridBuilder.NumWG()
}

// Next finds the location to dispatch the next work-group.
// TODO: WORKS ONLY FOR ONE APPLICATION AS OF NOW!!!
func (a *chipletLASPAlgorithm) Next() (location dispatchLocation) {
	a.numChiplets = 4

	cusPerChiplet := a.cuPool.NumCU() / a.numChiplets

	startingChiplet := a.nextChipletID
	for c := 0; c < a.numChiplets; c++ {
		chiplet := (startingChiplet + c) % a.numChiplets

		if len(a.perChipletListOfWG[chiplet]) != 0 {
			a.currWG = a.perChipletListOfWG[chiplet][0]
		} else {
			continue
		}

		if a.currWG == nil {
			continue
		}

		for i := 0; i < cusPerChiplet; i++ {
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

				a.currWG = nil
				a.perChipletListOfWG[chiplet] = a.perChipletListOfWG[chiplet][1:]
				a.numDispatchedWGs++

				return dispatch
			} else {
				//do nothing :)
			}
		}
	}
	return dispatchLocation{}
}

// FreeResources marks the dispatched location to be available.
func (a *chipletLASPAlgorithm) FreeResources(location dispatchLocation) {
	a.cuPool.GetCU(location.cuID).FreeResourcesForWG(location.wg)
}
