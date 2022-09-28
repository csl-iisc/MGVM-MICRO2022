package noc

import "gitlab.com/akita/akita"

// MsgBuffer is a buffer that can hold requests
type MsgBuffer struct {
	Capacity int
	Buf      []akita.Msg
	vc       int
}

func (b *MsgBuffer) enqueue(req akita.Msg) {
	if len(b.Buf) > b.Capacity {
		panic("buffer overflow")
	}

	b.Buf = append(b.Buf, req)
}
