package addresstranslator

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/device"
)

// A IdealAddressTranslatorBuilder can create address translators
type IdealAddressTranslatorBuilder struct {
	*Builder
	// engine          akita.Engine
	// freq            akita.Freq
	pageTable device.PageTable
	// ctrlPort        akita.Port
	// lowModuleFinder cache.LowModuleFinder
	// numReqPerCycle int
	// log2PageSize   uint64
	// deviceID       uint64
}

// MakeIdealAddressTranslatorBuilder creates a new IdealAddressTranslatorBuilder
func MakeIdealAddressTranslatorBuilder() AddressTranslatorBuilder {
	return &IdealAddressTranslatorBuilder{
		&Builder{
			freq:           1 * akita.GHz,
			numReqPerCycle: 4,
			log2PageSize:   12,
			deviceID:       1,
		}, nil}
}

// WithEngine sets the engine to be used by the address translators
// func (b IdealAddressTranslatorBuilder) WithEngine(engine akita.Engine) IdealAddressTranslatorBuilder {
// 	b.engine = engine
// 	return b
// }

// // WithFreq sets the frequency of the address translators
// func (b IdealAddressTranslatorBuilder) WithFreq(freq akita.Freq) IdealAddressTranslatorBuilder {
// 	b.freq = freq
// 	return b
// }

// WithTranslationProvider sets the port that can provide the translation
// service. The port must be a port on a TLB or an MMU.
// func (b IdealAddressTranslatorBuilder) WithTranslationProvider(p device.PageTable) {
// b.pageTable = p
// }

// WithLowModuleFinder sets the low modules finder that can tell the address
// translators where to send the memory access request to.
// func (b IdealAddressTranslatorBuilder) WithLowModuleFinder(f cache.LowModuleFinder) IdealAddressTranslatorBuilder {
// 	b.lowModuleFinder = f
// 	return b
// }

// WithNumReqPerCycle sets the number of request the address translators can
// process in each cycle.
// func (b IdealAddressTranslatorBuilder) WithNumReqPerCycle(n int) IdealAddressTranslatorBuilder {
// 	b.numReqPerCycle = n
// 	return b
// }

// // WithLog2PageSize sets the page size as a power of 2
// func (b IdealAddressTranslatorBuilder) WithLog2PageSize(n uint64) IdealAddressTranslatorBuilder {
// 	b.log2PageSize = n
// 	return b
// }

// // WithDeviceID sets the GPU ID that the address translator belongs to
// func (b IdealAddressTranslatorBuilder) WithDeviceID(n uint64) IdealAddressTranslatorBuilder {
// 	b.deviceID = n
// 	return b
// }

//WithCtrlPort sets the port of the component that can send ctrl reqs to AT
// func (b IdealAddressTranslatorBuilder) WithCtrlPort(p akita.Port) IdealAddressTranslatorBuilder {
// 	b.ctrlPort = p
// 	return b
// }

// Build returns a new AddressTranslator
func (b *IdealAddressTranslatorBuilder) Build(name string) AddressTranslator {
	t := &IdealAddressTranslator{}
	t.TickingComponent = akita.NewTickingComponent(
		name, b.engine, b.freq, t)

	t.TopPort = akita.NewLimitNumMsgPort(t, b.numReqPerCycle,
		name+".TopPort")
	t.BottomPort = akita.NewLimitNumMsgPort(t, b.numReqPerCycle,
		name+".BottomPort")
	t.CtrlPort = akita.NewLimitNumMsgPort(t, 1,
		name+".CtrlPort")

	t.pageTable = b.pageTable
	t.lowModuleFinder = b.lowModuleFinder
	t.numReqPerCycle = b.numReqPerCycle
	t.log2PageSize = b.log2PageSize
	t.deviceID = b.deviceID

	return t
}
