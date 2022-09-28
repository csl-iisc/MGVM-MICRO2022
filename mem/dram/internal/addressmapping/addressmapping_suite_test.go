package addressmapping

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAddressmapping(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Addressmapping Suite")
}
