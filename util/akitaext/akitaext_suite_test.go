package akitaext_test

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.com/akita/akita"
)

//go:generate mockgen -destination "mock_akita_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/akita Port,Engine,Ticker
//go:generate mockgen -destination "mock_util_test.go" -package $GOPACKAGE -write_package_comment=false gitlab.com/akita/util Buffer
func TestAkitaext(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Akitaext Suite")
}

type sampleMsg struct {
	akita.MsgMeta
}

func (m *sampleMsg) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}
