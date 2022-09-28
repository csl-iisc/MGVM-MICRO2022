package arbitration

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_util_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util Buffer
func TestArbitration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arbitration Suite")
}
