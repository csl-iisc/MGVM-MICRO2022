package idealmemcontroller

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port,Connection,Engine
func TestIdealmemcontroller(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Idealmemcontroller Suite")
}
