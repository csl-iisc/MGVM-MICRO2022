package cache

import (
	"log"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -destination "mock_cache_test.go" -package $GOPACKAGE  -write_package_comment=false -self_package=gitlab.com/akita/mem/cache gitlab.com/akita/mem/cache VictimFinder,Directory

func TestCache(t *testing.T) {
	log.SetOutput(GinkgoWriter)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cache Suite")
}
