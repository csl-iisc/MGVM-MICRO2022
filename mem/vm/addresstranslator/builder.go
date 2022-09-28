package addresstranslator

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem/cache"
)

type AddressTranslatorBuilder interface {
	Build(name string) AddressTranslator
	WithEngine(engine akita.Engine)
	WithFreq(freq akita.Freq)
	// WithTranslationProvider(p akita.Port)
	WithLowModuleFinder(f cache.LowModuleFinder)
	WithNumReqPerCycle(n int)
	WithLog2PageSize(n uint64)
	WithDeviceID(n uint64)
	WithCtrlPort(p akita.Port)
}

// A Builder can create address translators
type Builder struct {
	engine              akita.Engine
	freq                akita.Freq
	translationProvider akita.Port
	ctrlPort            akita.Port
	lowModuleFinder     cache.LowModuleFinder
	numReqPerCycle      int
	log2PageSize        uint64
	deviceID            uint64
}

// MakeBuilder creates a new builder
func MakeBuilder() AddressTranslatorBuilder {
	return &Builder{
		freq:           1 * akita.GHz,
		numReqPerCycle: 4,
		log2PageSize:   12,
		deviceID:       1,
	}
}

// WithEngine sets the engine to be used by the address translators
func (b *Builder) WithEngine(engine akita.Engine) {
	b.engine = engine
}

// WithFreq sets the frequency of the address translators
func (b *Builder) WithFreq(freq akita.Freq) {
	b.freq = freq
}

// WithTranslationProvider sets the port that can provide the translation
// service. The port must be a port on a TLB or an MMU.
// func (b *Builder) WithTranslationProvider(p akita.Port) {
// b.translationProvider = p

// }

// WithLowModuleFinder sets the low modules finder that can tell the address
// translators where to send the memory access request to.
func (b *Builder) WithLowModuleFinder(f cache.LowModuleFinder) {
	b.lowModuleFinder = f

}

// WithNumReqPerCycle sets the number of request the address translators can
// process in each cycle.
func (b *Builder) WithNumReqPerCycle(n int) {
	b.numReqPerCycle = n

}

// WithLog2PageSize sets the page size as a power of 2
func (b *Builder) WithLog2PageSize(n uint64) {
	b.log2PageSize = n

}

// WithDeviceID sets the GPU ID that the address translator belongs to
func (b *Builder) WithDeviceID(n uint64) {
	b.deviceID = n

}

//WithCtrlPort sets the port of the component that can send ctrl reqs to AT
func (b *Builder) WithCtrlPort(p akita.Port) {
	b.ctrlPort = p

}

//WithCtrlPort sets the port of the component that can send ctrl reqs to AT
// func (b *Builder) WithIdealVM() {
// 	b.useIdealVM = true

// }

// Build returns a new AddressTranslator
func (b *Builder) Build(name string) AddressTranslator {
	t := &DefaultAddressTranslator{}
	t.TickingComponent = akita.NewTickingComponent(
		name, b.engine, b.freq, t)

	t.TopPort = akita.NewLimitNumMsgPort(t, b.numReqPerCycle,
		name+".TopPort")
	t.BottomPort = akita.NewLimitNumMsgPort(t, b.numReqPerCycle,
		name+".BottomPort")
	t.TranslationPort = akita.NewLimitNumMsgPort(t, b.numReqPerCycle,
		name+".TranslationPort")
	t.CtrlPort = akita.NewLimitNumMsgPort(t, 1,
		name+".CtrlPort")

	t.translationProvider = b.translationProvider
	t.lowModuleFinder = b.lowModuleFinder
	t.numReqPerCycle = b.numReqPerCycle
	t.log2PageSize = b.log2PageSize
	t.deviceID = b.deviceID

	return t
}
