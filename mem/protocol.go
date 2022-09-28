package mem

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util/ca"
)

var accessReqByteOverhead = 12
var accessRspByteOverhead = 4
var controlMsgByteOverhead = 4

type AccessResult int

const (
	ReadHit = iota
	ReadMiss
	ReadMSHRHit
)

// AccessReq abstracts read and write requests that are sent to the
// cache modules or memory controllers.
type AccessReq interface {
	akita.Msg
	GetAddress() uint64
	GetByteSize() uint64
	GetPID() ca.PID
	// GetAccessInfo() interface{}
}

// A AccessRsp is a respond in the memory system.
type AccessRsp interface {
	akita.Msg
	GetRespondTo() string
}

// A ReadReq is a request sent to a memory controller to fetch data
type ReadReq struct {
	akita.MsgMeta

	Address            uint64
	AccessByteSize     uint64
	PID                ca.PID
	CanWaitForCoalesce bool
	Info               interface{}
}

// Meta returns the message meta.
func (r *ReadReq) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// GetByteSize returns the number of byte that the request is accessing.
func (r *ReadReq) GetByteSize() uint64 {
	return r.AccessByteSize
}

// GetAddress returns the address that the request is accessing
func (r *ReadReq) GetAddress() uint64 {
	return r.Address
}

// GetPID returns the process ID that the request is working on.
func (r *ReadReq) GetPID() ca.PID {
	return r.PID
}

// GetAccessInfo returns the process ID that the request is working on.
// func (r *ReadReq) GetAccessInfo() interface{} {
// 	return r.Info
// }

// A ReadReqInfo is a request sent to a memory controller to write data
type ReadReqInfo struct {
	ReturnAccessInfo bool
	AccessResult
}

// ReadReqBuilder can build read requests.
type ReadReqBuilder struct {
	sendTime           akita.VTimeInSec
	src, dst           akita.Port
	pid                ca.PID
	address, byteSize  uint64
	canWaitForCoalesce bool
	info               interface{}
}

// WithSendTime sets the send time of the request to build.
func (b ReadReqBuilder) WithSendTime(t akita.VTimeInSec) ReadReqBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b ReadReqBuilder) WithSrc(src akita.Port) ReadReqBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b ReadReqBuilder) WithDst(dst akita.Port) ReadReqBuilder {
	b.dst = dst
	return b
}

// WithPID sets the PID of the request to build.
func (b ReadReqBuilder) WithPID(pid ca.PID) ReadReqBuilder {
	b.pid = pid
	return b
}

// WithInfo sets the Info of the request to build.
func (b ReadReqBuilder) WithInfo(info interface{}) ReadReqBuilder {
	b.info = info
	return b
}

// WithAddress sets the address of the request to build.
func (b ReadReqBuilder) WithAddress(address uint64) ReadReqBuilder {
	b.address = address
	return b
}

// WithByteSize sets the byte size of the request to build.
func (b ReadReqBuilder) WithByteSize(byteSize uint64) ReadReqBuilder {
	b.byteSize = byteSize
	return b
}

// CanWaitForCoalesce allow the request to build to wait for coalesce.
func (b ReadReqBuilder) CanWaitForCoalesce() ReadReqBuilder {
	b.canWaitForCoalesce = true
	return b
}

// Build creates a new ReadReq
func (b ReadReqBuilder) Build() *ReadReq {
	r := &ReadReq{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.TrafficBytes = accessReqByteOverhead
	r.Address = b.address
	r.PID = b.pid
	r.Info = b.info
	r.AccessByteSize = b.byteSize
	r.CanWaitForCoalesce = b.canWaitForCoalesce
	return r
}

// A WriteReq is a request sent to a memory controller to write data
type WriteReq struct {
	akita.MsgMeta

	Address            uint64
	Data               []byte
	DirtyMask          []bool
	PID                ca.PID
	CanWaitForCoalesce bool
	Info               interface{}
}

// Meta returns the meta data attached to a request.
func (r *WriteReq) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// GetByteSize returns the number of byte that the request is writing.
func (r *WriteReq) GetByteSize() uint64 {
	return uint64(len(r.Data))
}

// GetAddress returns the address that the request is accessing
func (r *WriteReq) GetAddress() uint64 {
	return r.Address
}

// GetPID returns the PID of the read address
func (r *WriteReq) GetPID() ca.PID {
	return r.PID
}

// GetAccessInfo returns the process ID that the request is working on.
// func (r *WriteReq) GetAccessInfo() interface{} {
// 	return r.Info
// }

// A WriteReqInfo is a request sent to a memory controller to write data
type WriteReqInfo struct {
	returnAccessInfo bool
}

// WriteReqBuilder can build read requests.
type WriteReqBuilder struct {
	sendTime           akita.VTimeInSec
	src, dst           akita.Port
	pid                ca.PID
	info               interface{}
	address            uint64
	data               []byte
	dirtyMask          []bool
	canWaitForCoalesce bool
}

// WithSendTime sets the send time of the message to build.
func (b WriteReqBuilder) WithSendTime(t akita.VTimeInSec) WriteReqBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b WriteReqBuilder) WithSrc(src akita.Port) WriteReqBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b WriteReqBuilder) WithDst(dst akita.Port) WriteReqBuilder {
	b.dst = dst
	return b
}

// WithPID sets the PID of the request to build.
func (b WriteReqBuilder) WithPID(pid ca.PID) WriteReqBuilder {
	b.pid = pid
	return b
}

// WithInfo sets the information attached to the request to build.
func (b WriteReqBuilder) WithInfo(info interface{}) WriteReqBuilder {
	b.info = info
	return b
}

// WithAddress sets the address of the request to build.
func (b WriteReqBuilder) WithAddress(address uint64) WriteReqBuilder {
	b.address = address
	return b
}

// WithData sets the data of the request to build.
func (b WriteReqBuilder) WithData(data []byte) WriteReqBuilder {
	b.data = data
	return b
}

// WithDirtyMask sets the dirty mask of the request to build.
func (b WriteReqBuilder) WithDirtyMask(mask []bool) WriteReqBuilder {
	b.dirtyMask = mask
	return b
}

// CanWaitForCoalesce allow the request to build to wait for coalesce.
func (b WriteReqBuilder) CanWaitForCoalesce() WriteReqBuilder {
	b.canWaitForCoalesce = true
	return b
}

// Build creates a new WriteReq
func (b WriteReqBuilder) Build() *WriteReq {
	r := &WriteReq{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.PID = b.pid
	r.Info = b.info
	r.Address = b.address
	r.Data = b.data
	r.TrafficBytes = len(r.Data) + accessReqByteOverhead
	r.DirtyMask = b.dirtyMask
	r.CanWaitForCoalesce = b.canWaitForCoalesce
	return r
}

// A DataReadyRsp is the respond sent from the lower module to the higher
// module that carries the data loaded.
type DataReadyRsp struct {
	akita.MsgMeta

	RespondTo string // The ID of the request it replies
	Data      []byte
	Info      interface{}
}

// Meta returns the meta data attached to each message.
func (r *DataReadyRsp) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// GetRespondTo returns the ID if the request that the respond is resonding to.
func (r *DataReadyRsp) GetRespondTo() string {
	return r.RespondTo
}

// DataReadyRspInfo stores
type DataReadyRspInfo struct {
	AccessResult
	Src string
}

// DataReadyRspBuilder can build data ready responds.
type DataReadyRspBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
	rspTo    string
	data     []byte
	info     interface{}
}

// WithSendTime sets the send time of the request to build.
func (b DataReadyRspBuilder) WithSendTime(
	t akita.VTimeInSec,
) DataReadyRspBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b DataReadyRspBuilder) WithSrc(src akita.Port) DataReadyRspBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b DataReadyRspBuilder) WithDst(dst akita.Port) DataReadyRspBuilder {
	b.dst = dst
	return b
}

// WithRspTo sets ID of the request that the respond to build is replying to.
func (b DataReadyRspBuilder) WithRspTo(id string) DataReadyRspBuilder {
	b.rspTo = id
	return b
}

// WithData sets the data of the request to build.
func (b DataReadyRspBuilder) WithData(data []byte) DataReadyRspBuilder {
	b.data = data
	return b
}

// WithInfo sets the info of the request to build.
func (b DataReadyRspBuilder) WithInfo(info interface{}) DataReadyRspBuilder {
	b.info = info
	return b
}

// Build creates a new DataReadyRsp
func (b DataReadyRspBuilder) Build() *DataReadyRsp {
	r := &DataReadyRsp{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.TrafficBytes = len(b.data) + accessRspByteOverhead
	r.RespondTo = b.rspTo
	r.Data = b.data
	r.Info = b.info
	return r
}

// A WriteDoneRsp is a respond sent from the lower module to the higher module
// to mark a previous requests is completed successfully.
type WriteDoneRsp struct {
	akita.MsgMeta

	RespondTo string
}

// Meta returns the meta data accociated with the message.
func (r *WriteDoneRsp) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// GetRespondTo returns the ID of the request that the respond is responding to.
func (r *WriteDoneRsp) GetRespondTo() string {
	return r.RespondTo
}

// WriteDoneRspBuilder can build data ready responds.
type WriteDoneRspBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
	rspTo    string
}

// WithSendTime sets the send time of the message to build.
func (b WriteDoneRspBuilder) WithSendTime(
	t akita.VTimeInSec,
) WriteDoneRspBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b WriteDoneRspBuilder) WithSrc(src akita.Port) WriteDoneRspBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b WriteDoneRspBuilder) WithDst(dst akita.Port) WriteDoneRspBuilder {
	b.dst = dst
	return b
}

// WithRspTo sets ID of the request that the respond to build is replying to.
func (b WriteDoneRspBuilder) WithRspTo(id string) WriteDoneRspBuilder {
	b.rspTo = id
	return b
}

// Build creates a new WriteDoneRsp
func (b WriteDoneRspBuilder) Build() *WriteDoneRsp {
	r := &WriteDoneRsp{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.TrafficBytes = accessRspByteOverhead
	r.SendTime = b.sendTime
	r.RespondTo = b.rspTo
	return r
}

// ControlMsg is the commonly used message type for controlling the components
// on the memory hierarchy. It is also used for resonpding the original
// requester with the Done field.
type ControlMsg struct {
	akita.MsgMeta

	DiscardTransations bool
	Restart            bool
	NotifyDone         bool
}

// Meta returns the meta data assocated with the ControlMsg.
func (m *ControlMsg) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}

// A ControlMsgBuilder can build control messages.
type ControlMsgBuilder struct {
	sendTime            akita.VTimeInSec
	src, dst            akita.Port
	discardTransactions bool
	restart             bool
	notifyDone          bool
}

// WithSendTime sets the send time of the message to build.
func (b ControlMsgBuilder) WithSendTime(
	t akita.VTimeInSec,
) ControlMsgBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b ControlMsgBuilder) WithSrc(src akita.Port) ControlMsgBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b ControlMsgBuilder) WithDst(dst akita.Port) ControlMsgBuilder {
	b.dst = dst
	return b
}

// ToDiscardTransactions sets the discard transactions bit of the control
// messages to 1.
func (b ControlMsgBuilder) ToDiscardTransactions() ControlMsgBuilder {
	b.discardTransactions = true
	return b
}

// ToRestart sets the restart bit of the control messages to 1.
func (b ControlMsgBuilder) ToRestart() ControlMsgBuilder {
	b.restart = true
	return b
}

// ToNotifyDone sets the "notify done" bit of the control messages to 1.
func (b ControlMsgBuilder) ToNotifyDone() ControlMsgBuilder {
	b.notifyDone = true
	return b
}

func (b ControlMsgBuilder) Build() *ControlMsg {
	m := &ControlMsg{}
	m.ID = akita.GetIDGenerator().Generate()
	m.Src = b.src
	m.Dst = b.dst
	m.TrafficBytes = controlMsgByteOverhead
	m.SendTime = b.sendTime

	m.DiscardTransations = b.discardTransactions
	m.Restart = b.restart
	m.NotifyDone = b.notifyDone

	return m
}
