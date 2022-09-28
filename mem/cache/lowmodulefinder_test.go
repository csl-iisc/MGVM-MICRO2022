package cache

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
)

var _ = Describe("InterleavedLowModuleFinder", func() {
	var (
		lowModuleFinder *InterleavedLowModuleFinder
	)

	BeforeEach(func() {
		lowModuleFinder = new(InterleavedLowModuleFinder)
		lowModuleFinder.UseAddressSpaceLimitation = true
		lowModuleFinder.LowAddress = 0
		lowModuleFinder.HighAddress = 4 * mem.GB
		lowModuleFinder.InterleavingSize = 4096
		lowModuleFinder.LowModules = make([]akita.Port, 0)
		for i := 0; i < 6; i++ {
			lowModuleFinder.LowModules = append(
				lowModuleFinder.LowModules,
				akita.NewLimitNumMsgPort(nil, 4,
					fmt.Sprintf("LowModule_%d.Port", i)))
		}
		lowModuleFinder.ModuleForOtherAddresses =
			akita.NewLimitNumMsgPort(nil, 4, "LowMoudle_other.Port")
	})

	It("should find low module if address is in-space", func() {
		Expect(lowModuleFinder.Find(0)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[0]))
		Expect(lowModuleFinder.Find(4096)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[1]))
		Expect(lowModuleFinder.Find(4097)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[1]))
	})

	It("should use a special module for all the addresses that does not fall in range", func() {
		Expect(lowModuleFinder.Find(4 * mem.GB)).To(
			BeIdenticalTo(lowModuleFinder.ModuleForOtherAddresses))
	})
})
