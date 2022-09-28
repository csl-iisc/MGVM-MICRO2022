package cache

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
)

var _ = Describe("TwoLevelLowModuleFinder", func() {
	var (
		lowModuleFinder *TwoLevelLowModuleFinder
	)

	BeforeEach(func() {
		lowModuleFinder = NewTwoLevelLowModuleFinder(
			4*mem.GB, 8*mem.GB,
			2*mem.GB, 2,
			4096, 4)
		for i := 0; i < 8; i++ {
			lowModuleFinder.LowModules = append(
				lowModuleFinder.LowModules,
				akita.NewLimitNumMsgPort(nil, 4,
					fmt.Sprintf("LowModule_%d.Port", i)))
		}
	})

	It("should find out the module based on the above TwoLevel", func() {
		Expect(lowModuleFinder.Find(4295507968)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[0]))
		Expect(lowModuleFinder.Find(4295512064)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[1]))
		Expect(lowModuleFinder.Find(4295512085)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[1]))
		Expect(lowModuleFinder.Find(6442450944)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[4]))
		Expect(lowModuleFinder.Find(6442459136)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[6]))
	})

})

var _ = Describe("LocalInterleavedLowModuleFinder", func() {
	var (
		lowModuleFinder *LocalInterleavedLowModuleFinder
	)

	BeforeEach(func() {
		lowModuleFinder = NewLocalInterleavedLowModuleFinder(2, 4, 32)
		lowModuleFinder.LocalLowModule = akita.NewLimitNumMsgPort(nil, 4, "local")
		lowModuleFinder.RemoteLowModule = akita.NewLimitNumMsgPort(nil, 4, "remote")
	})

	It("should find out the module based on the above LocalInterleaved", func() {
		Expect(lowModuleFinder.Find(4)).To(
			BeIdenticalTo(lowModuleFinder.RemoteLowModule))
		Expect(lowModuleFinder.Find(66)).To(
			BeIdenticalTo(lowModuleFinder.LocalLowModule))
	})

})
