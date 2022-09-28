package arbitration

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
)

// Arbiter can determine which buffer can send a message out
type Arbiter interface {
	// Add a buffer for arbitration
	AddBuffer(buf util.Buffer)

	// Arbitrate returns a set of ports that can send request in the next cycle.
	Arbitrate(now akita.VTimeInSec) []util.Buffer
}
