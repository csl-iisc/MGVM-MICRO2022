package cache

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
)

var _ = Describe("HashLowModuleFinder", func() {
	var (
		lowModuleFinder *HashLowModuleFinder
	)

	BeforeEach(func() {
		lowModuleFinder = NewHashLowModuleFinder()
		for i := 0; i < 7; i++ {
			lowModuleFinder.LowModules = append(
				lowModuleFinder.LowModules,
				akita.NewLimitNumMsgPort(nil, 4,
					fmt.Sprintf("LowModule_%d.Port", i)))
		}
	})

	It("should find out the module based on the above hash", func() {
		Expect(lowModuleFinder.Find(0)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[0]))
		Expect(lowModuleFinder.Find(43)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[1]))
		Expect(lowModuleFinder.Find(23)).To(
			BeIdenticalTo(lowModuleFinder.LowModules[2]))
	})

})
