package addresstranslator

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port,Engine
//go:generate mockgen -destination "mock_cache_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/cache LowModuleFinder
//go:generate mockgen -destination "mock_akitaext_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util/akitaext BufferedSender
func TestAddresstranslator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Address Translator Suite")
}
