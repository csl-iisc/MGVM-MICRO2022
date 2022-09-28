package tlb

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util/ca"
)

// A TLBFlushReq asks the TLB to invalidate certain entries. It will also not block all incoming and outgoing ports
type TLBFlushReq struct {
	akita.MsgMeta
	VAddr []uint64
	PID   ca.PID
}

// Meta returns the meta data associated with the message.
func (r *TLBFlushReq) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TLBFlushReqBuilder can build AT flush requests
type TLBFlushReqBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
	vAddrs   []uint64
	pid      ca.PID
}

// WithSendTime sets the send time of the request to build.:w
func (b TLBFlushReqBuilder) WithSendTime(
	t akita.VTimeInSec,
) TLBFlushReqBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b TLBFlushReqBuilder) WithSrc(src akita.Port) TLBFlushReqBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b TLBFlushReqBuilder) WithDst(dst akita.Port) TLBFlushReqBuilder {
	b.dst = dst
	return b
}

// WithVAddrs sets the Vaddr of the pages to be flushed
func (b TLBFlushReqBuilder) WithVAddrs(vAddrs []uint64) TLBFlushReqBuilder {
	b.vAddrs = vAddrs
	return b
}

// WithPID sets the pid whose entries are to be flushed
func (b TLBFlushReqBuilder) WithPID(pid ca.PID) TLBFlushReqBuilder {
	b.pid = pid
	return b
}

// Build creates a new TLBFlushReq
func (b TLBFlushReqBuilder) Build() *TLBFlushReq {
	r := &TLBFlushReq{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.VAddr = b.vAddrs
	r.PID = b.pid
	return r
}

// A TLBFlushRsp is a response from AT indicating flush is complete
type TLBFlushRsp struct {
	akita.MsgMeta
}

// Meta returns the meta data associated with the message.
func (r *TLBFlushRsp) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TLBFlushRspBuilder can build AT flush rsp
type TLBFlushRspBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
}

// WithSendTime sets the send time of the request to build.:w
func (b TLBFlushRspBuilder) WithSendTime(
	t akita.VTimeInSec,
) TLBFlushRspBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b TLBFlushRspBuilder) WithSrc(src akita.Port) TLBFlushRspBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b TLBFlushRspBuilder) WithDst(dst akita.Port) TLBFlushRspBuilder {
	b.dst = dst
	return b
}

// Build creates a new TLBFlushRsps.
func (b TLBFlushRspBuilder) Build() *TLBFlushRsp {
	r := &TLBFlushRsp{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime

	return r
}

// A TLBRestartReq is a request to TLB to start accepting requests and resume operations
type TLBRestartReq struct {
	akita.MsgMeta
}

// Meta returns the meta data associated with the message.
func (r *TLBRestartReq) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TLBRestartReqBuilder can build TLB restart requests.
type TLBRestartReqBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
}

// WithSendTime sets the send time of the request to build.
func (b TLBRestartReqBuilder) WithSendTime(
	t akita.VTimeInSec,
) TLBRestartReqBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b TLBRestartReqBuilder) WithSrc(src akita.Port) TLBRestartReqBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b TLBRestartReqBuilder) WithDst(dst akita.Port) TLBRestartReqBuilder {
	b.dst = dst
	return b
}

// Build creates a new TLBRestartReq.
func (b TLBRestartReqBuilder) Build() *TLBRestartReq {
	r := &TLBRestartReq{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime

	return r
}

// A TLBRestartRsp is a response from AT indicating it has resumed working
type TLBRestartRsp struct {
	akita.MsgMeta
}

// Meta returns the meta data associated with the message.
func (r *TLBRestartRsp) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TLBRestartRspBuilder can build AT flush rsp
type TLBRestartRspBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
}

// WithSendTime sets the send time of the request to build.:w
func (b TLBRestartRspBuilder) WithSendTime(
	t akita.VTimeInSec,
) TLBRestartRspBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b TLBRestartRspBuilder) WithSrc(src akita.Port) TLBRestartRspBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b TLBRestartRspBuilder) WithDst(dst akita.Port) TLBRestartRspBuilder {
	b.dst = dst
	return b
}

// Build creates a new TLBRestartRsp
func (b TLBRestartRspBuilder) Build() *TLBRestartRsp {
	r := &TLBRestartRsp{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime

	return r
}
