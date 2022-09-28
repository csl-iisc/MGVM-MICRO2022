package noc

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util"
)

// Flit is the smallest trasferring unit on a network.
type Flit struct {
	akita.MsgMeta
	SeqID        int
	NumFlitInMsg int
	Msg          akita.Msg
	OutputBuf    util.Buffer // The buffer to route to within a switch
}

// Meta returns the meta data assocated with the Flit.
func (f *Flit) Meta() *akita.MsgMeta {
	return &f.MsgMeta
}

// FlitBuilder can build flits
type FlitBuilder struct {
	sendTime            akita.VTimeInSec
	src, dst            akita.Port
	msg                 akita.Msg
	seqID, numFlitInMsg int
}

// WithSendTime sets the send time of the request to build
func (b FlitBuilder) WithSendTime(t akita.VTimeInSec) FlitBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the src of the request to send
func (b FlitBuilder) WithSrc(src akita.Port) FlitBuilder {
	b.src = src
	return b
}

// WithDst sets the dst of the request to send
func (b FlitBuilder) WithDst(dst akita.Port) FlitBuilder {
	b.dst = dst
	return b
}

// WithSeqID sets the SeqID of the Flit.
func (b FlitBuilder) WithSeqID(i int) FlitBuilder {
	b.seqID = i
	return b
}

// WithNumFlitInMsg sets the NumFlitInMsg for of flit to build.
func (b FlitBuilder) WithNumFlitInMsg(n int) FlitBuilder {
	b.numFlitInMsg = n
	return b
}

// WithMsg sets the msg of the flit to build.
func (b FlitBuilder) WithMsg(msg akita.Msg) FlitBuilder {
	b.msg = msg
	return b
}

// Build creates a new flit.
func (b FlitBuilder) Build() *Flit {
	f := &Flit{}
	f.SendTime = b.sendTime
	f.Src = b.src
	f.Dst = b.dst
	f.Msg = b.msg
	f.SeqID = b.seqID
	f.NumFlitInMsg = b.numFlitInMsg
	return f
}
