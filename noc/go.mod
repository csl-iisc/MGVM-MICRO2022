module gitlab.com/akita/noc

require (
	github.com/golang/mock v1.4.3
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/tebeka/atexit v0.3.0
	gitlab.com/akita/akita v1.10.1
	gitlab.com/akita/util v0.4.0
	google.golang.org/appengine v1.6.5 // indirect
)

// replace gitlab.com/akita/akita => ../akita

// replace gitlab.com/akita/util => ../util

go 1.13
