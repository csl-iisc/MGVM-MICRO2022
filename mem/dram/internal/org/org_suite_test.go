package org

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -source bank.go -destination mock_bank_test.go -self_package gitlab.com/akita/mem/dram/internal/org -package $GOPACKAGE

func TestOrg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Org Suite")
}
