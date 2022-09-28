package noc

import (
	"gitlab.com/akita/akita"
)

// A TrafficCounter counts number of bytes transferred over a connection
type TrafficCounter struct {
	TotalData uint64
}

// Func adds the delivered traffic to the counter
func (c *TrafficCounter) Func(ctx *akita.HookCtx) {
	if ctx.Pos != akita.HookPosConnDeliver {
		return
	}

	req := ctx.Item.(akita.Msg)
	c.TotalData += uint64(req.Meta().TrafficBytes)
}
