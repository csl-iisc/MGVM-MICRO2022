// Package vm provides the models for address translations
package device

import (
	"gitlab.com/akita/akita"
	"gitlab.com/akita/util/ca"
)

type AccessResult int

const (
	TLBUnitialized = iota
	TLBHit
	TLBMiss
	TLBMshrHit
)

// A TranslationReq asks the receiver component to translate the request.
type TranslationReq struct {
	akita.MsgMeta
	VAddr    uint64
	PID      ca.PID
	DeviceID uint64
}

// Meta returns the meta data associated with the message.
func (r *TranslationReq) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TranslationReqBuilder can build translation requests
type TranslationReqBuilder struct {
	sendTime akita.VTimeInSec
	src, dst akita.Port
	vAddr    uint64
	pid      ca.PID
	deviceID uint64
}

// WithSendTime sets the send time of the request to build.:w
func (b TranslationReqBuilder) WithSendTime(
	t akita.VTimeInSec,
) TranslationReqBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the request to build.
func (b TranslationReqBuilder) WithSrc(src akita.Port) TranslationReqBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the request to build.
func (b TranslationReqBuilder) WithDst(dst akita.Port) TranslationReqBuilder {
	b.dst = dst
	return b
}

// WithVAddr sets the virtual address of the request to build.
func (b TranslationReqBuilder) WithVAddr(vAddr uint64) TranslationReqBuilder {
	b.vAddr = vAddr
	return b
}

// WithPID sets the virtual address of the request to build.
func (b TranslationReqBuilder) WithPID(pid ca.PID) TranslationReqBuilder {
	b.pid = pid
	return b
}

// WithdeviceID sets the GPU ID of the request to build.
func (b TranslationReqBuilder) WithDeviceID(deviceID uint64) TranslationReqBuilder {
	b.deviceID = deviceID
	return b
}

// Build creats a new TranslationReq
func (b TranslationReqBuilder) Build() *TranslationReq {
	r := &TranslationReq{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.VAddr = b.vAddr
	r.PID = b.pid
	r.DeviceID = b.deviceID
	return r
}

// A TranslationRsp is the respond for a TranslationReq. It carries the physical
// address.
type TranslationRsp struct {
	akita.MsgMeta
	RespondTo string // The ID of the request it replies
	Page      Page
	HitOrMiss AccessResult
	SrcL2TLB  string
}

// Meta returns the meta data associated with the message.
func (r *TranslationRsp) Meta() *akita.MsgMeta {
	return &r.MsgMeta
}

// TranslationRspBuilder can build translation requests
type TranslationRspBuilder struct {
	sendTime  akita.VTimeInSec
	src, dst  akita.Port
	rspTo     string
	page      Page
	hitOrMiss AccessResult
	srcL2TLB  string
}

// WithSendTime sets the send time of the message to build.
func (b TranslationRspBuilder) WithSendTime(
	t akita.VTimeInSec,
) TranslationRspBuilder {
	b.sendTime = t
	return b
}

// WithSrc sets the source of the respond to build.
func (b TranslationRspBuilder) WithSrc(src akita.Port) TranslationRspBuilder {
	b.src = src
	return b
}

// WithDst sets the destination of the respond to build.
func (b TranslationRspBuilder) WithDst(dst akita.Port) TranslationRspBuilder {
	b.dst = dst
	return b
}

// WithRspTo sets the request ID of the respond to build.
func (b TranslationRspBuilder) WithRspTo(rspTo string) TranslationRspBuilder {
	b.rspTo = rspTo
	return b
}

// WithPage sets the page of the respond to build.
func (b TranslationRspBuilder) WithPage(page Page) TranslationRspBuilder {
	b.page = page
	return b
}

// WithAccessType sets the AccessType of the respond to build.
func (b TranslationRspBuilder) WithAccessResult(hitOrMiss AccessResult) TranslationRspBuilder {
	b.hitOrMiss = hitOrMiss
	return b
}

//WithAccessType sets the AccessType of the respond to build.
func (b TranslationRspBuilder) WithSrcL2TLB(srcL2TLB string) TranslationRspBuilder {
	b.srcL2TLB = srcL2TLB
	return b
}

// Build creats a new TranslationRsp
func (b TranslationRspBuilder) Build() *TranslationRsp {
	r := &TranslationRsp{}
	r.ID = akita.GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.RespondTo = b.rspTo
	r.Page = b.page
	r.HitOrMiss = b.hitOrMiss
	r.SrcL2TLB = b.srcL2TLB
	return r
}

type PageMigrationInfo struct {
	GPUReqToVAddrMap map[uint64][]uint64
}

//PageMigrationReqToDriver is a req to driver from MMU to start page migration process
type PageMigrationReqToDriver struct {
	akita.MsgMeta

	StartTime         akita.VTimeInSec
	EndTime           akita.VTimeInSec
	MigrationInfo     *PageMigrationInfo
	CurrAccessingGPUs []uint64
	PID               ca.PID
	CurrPageHostGPU   uint64
	PageSize          uint64
	RespondToTop      bool
}

// Meta returns the meta data associated with the message.
func (m *PageMigrationReqToDriver) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}

func NewPageMigrationReqToDriver(
	time akita.VTimeInSec,
	src, dst akita.Port,
) *PageMigrationReqToDriver {
	cmd := new(PageMigrationReqToDriver)
	cmd.SendTime = time
	cmd.Src = src
	cmd.Dst = dst
	return cmd
}

//PageMigrationRspFromDriver is a rsp from driver to MMU marking completion of migration
type PageMigrationRspFromDriver struct {
	akita.MsgMeta

	StartTime akita.VTimeInSec
	EndTime   akita.VTimeInSec
	VAddr     []uint64
	RspToTop  bool
}

// Meta returns the meta data associated with the message.
func (m *PageMigrationRspFromDriver) Meta() *akita.MsgMeta {
	return &m.MsgMeta
}

func NewPageMigrationRspFromDriver(
	time akita.VTimeInSec,
	src, dst akita.Port,
) *PageMigrationRspFromDriver {
	cmd := new(PageMigrationRspFromDriver)
	cmd.SendTime = time
	cmd.Src = src
	cmd.Dst = dst
	return cmd
}
