package cmdq

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_org_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/mem/dram/internal/org Channel

func TestCmdq(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmdq Suite")
}
