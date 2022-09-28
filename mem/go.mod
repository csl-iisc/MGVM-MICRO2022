module gitlab.com/akita/mem

require (
	github.com/golang/mock v1.4.4
	github.com/google/btree v1.0.0
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/rs/xid v1.2.1
	gitlab.com/akita/akita v1.10.1
	gitlab.com/akita/util v0.6.1
)

replace gitlab.com/akita/akita => ../akita

replace gitlab.com/akita/util => ../util

go 1.13
