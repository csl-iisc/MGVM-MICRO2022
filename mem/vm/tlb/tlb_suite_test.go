package tlb

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port,Engine
//go:generate mockgen -destination "mock_cache_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/cache LowModuleFinder
//go:generate mockgen -destination "mock_akitaext_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util/akitaext BufferedSender
//go:generate mockgen -destination "mock_internal_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/vm/tlb/internal Set
func TestTlb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tlb Suite")
}
