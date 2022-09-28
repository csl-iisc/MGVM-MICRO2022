package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mgpusim"
)

type DistributedTLBGPUBuilder struct {
	*CommonBuilder

	// specific componenets
}

// Distributed TLB specific function

// MakeDistributedTLBGPUBuilder provides a GPU builder that can builds MCM GPU.
func MakeDistributedTLBGPUBuilder() DistributedTLBGPUBuilder {
	// TODO: should this be using new? is the object being allocated on the stack?
	cbp := CommonBuilder{}
	b := DistributedTLBGPUBuilder{CommonBuilder: &cbp}
	b.SetDefaultCommonBuilderParams()
	return b
}

func (b DistributedTLBGPUBuilder) Build(name string, id uint64) *mgpusim.GPU {
	b.createGPU(name, id)

	b.buildCP()

	chipRdmaAddressTable := b.createChipRDMAAddrTable()
	rdmaResponsePorts := make([]akita.Port, 4)
	remoteAddressTranslationTable := b.createRemoteAddrTransTable()
	rtuResponsePorts := make([]akita.Port, 4)

	for i := 0; i < b.numChiplet; i++ {
		chipletName := fmt.Sprintf("%s.chiplet_%02d", b.gpuName, i)
		chiplet := NewChiplet(chipletName, uint64(i))

		b.BuildSAs(chiplet)
		b.buildMemBanks(chiplet)
		b.buildMMU(chiplet)
		b.buildL2TLB(chiplet)

		b.configChipRDMAEngine(chiplet, chipRdmaAddressTable, rdmaResponsePorts)
		b.configRemoteAddressTranslationUnit(chiplet, remoteAddressTranslationTable, rtuResponsePorts)

		b.connectL1ToL2(chiplet)
		b.connectL2ToDRAM(chiplet)
		b.connectL1TLBToL2TLB(chiplet)
		b.connectL2TLBTOMMU(chiplet)
		b.connectMMUToL2(chiplet)

		b.chiplets = append(b.chiplets, chiplet)
	}

	b.buildPageMigrationController()
	b.setupDMA()

	b.connectCP()
	b.setupInterchipNetwork()

	return b.gpu
}
