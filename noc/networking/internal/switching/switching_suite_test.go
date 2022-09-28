package switching

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port,Engine
//go:generate mockgen -destination "mock_util_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util Buffer
//go:generate mockgen -destination "mock_pipelining_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util/pipelining Pipeline
//go:generate mockgen -destination "mock_routing_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/noc/networking/internal/routing Table
//go:generate mockgen -destination "mock_arbitration_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/noc/networking/internal/arbitration Arbiter

func TestSwitching(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Switching Suite")
}

type sampleMsg struct {
	akita.MsgMeta
}

func (m *sampleMsg) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}
