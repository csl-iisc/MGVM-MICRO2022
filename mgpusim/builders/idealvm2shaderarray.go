package builders

import (
	"fmt"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
	"gitlab.com/akita/mem/vm/addresstranslator"
)

type idealVM2ShaderArray struct {
	*shaderArray
}

type idealVM2ShaderArrayBuilder struct {
	shaderArrayBuilder
	pageTable device.PageTable
}

func makeIdealVM2ShaderArrayBuilder() idealVM2ShaderArrayBuilder {
	b := idealVM2ShaderArrayBuilder{shaderArrayBuilder{
		gpuID:             0,
		name:              "SA",
		numCU:             4,
		freq:              1 * akita.GHz,
		log2CacheLineSize: 6,
		log2PageSize:      12,
	}, nil}
	return b
}

func (b *idealVM2ShaderArrayBuilder) withPageTable(
	pageTable device.PageTable,
) {
	b.pageTable = pageTable
}

func (b *idealVM2ShaderArrayBuilder) makeAddressTranslatorBuilder() (builder addresstranslator.AddressTranslatorBuilder) {
	builder = addresstranslator.MakeIdealAddressTranslatorBuilder()
	builder.WithEngine(b.engine)
	builder.WithFreq(b.freq)
	builder.WithDeviceID(b.gpuID)
	builder.WithLog2PageSize(b.log2PageSize)
	return
}
func (b *idealVM2ShaderArrayBuilder) makeIdealAddressTranslator(name string) addresstranslator.AddressTranslator {
	builder := b.makeAddressTranslatorBuilder()
	at := builder.Build(name)
	at.SetTranslationProvider(b.pageTable)
	return at
}

func (b *idealVM2ShaderArrayBuilder) buildAddressTranslators(sa *shaderArray) {
	for i := 0; i < b.numCU; i++ {
		name := fmt.Sprintf("%s.L1VAddrTrans_%02d", b.name, i)
		at := b.makeIdealAddressTranslator(name)
		sa.l1vATs = append(sa.l1vATs, at)
	}
	name := fmt.Sprintf("%s.L1SAddrTrans", b.name)
	at := b.makeIdealAddressTranslator(name)
	sa.l1sAT = at

	name = fmt.Sprintf("%s.L1IAddrTrans", b.name)
	at = b.makeIdealAddressTranslator(name)
	sa.l1iAT = at
}

func (b *idealVM2ShaderArrayBuilder) buildComponents(sa *shaderArray) {
	b.buildCUs(sa)

	b.buildAddressTranslators(sa)

	b.buildL1VReorderBuffers(sa)
	b.buildL1VCaches(sa)

	b.buildL1SReorderBuffer(sa)
	b.buildL1SCache(sa)

	b.buildL1IReorderBuffer(sa)
	b.buildL1ICache(sa)
}

func (b *idealVM2ShaderArrayBuilder) connectComponents(sa *shaderArray) {
	b.connectVectorMem(sa)
	b.connectScalarMem(sa)
	b.connectInstMem(sa)
}

func (b *idealVM2ShaderArrayBuilder) Build(name string) shaderArray {
	b.name = name
	sa := shaderArray{}

	b.buildComponents(&sa)
	b.connectComponents(&sa)

	return sa
}
