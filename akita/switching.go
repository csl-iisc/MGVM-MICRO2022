package akita

// We can move this to vm/tlb/tlbprotocol.go if necessay.

type TLBIndexingSwitch int

const (
	TLBIndexingSwitch4K TLBIndexingSwitch = iota
	TLBIndexingSwitchHSL
)

// TLBIndexingSwitchMsg is sent from CP to RTUs/TLBs requiring them to switch
// to use a different HSL
type TLBIndexingSwitchMsg struct {
	MsgMeta
	TLBInterleaving   uint64
	TLBIndexingSwitch TLBIndexingSwitch
}

func (t *TLBIndexingSwitchMsg) Meta() *MsgMeta {
	return &t.MsgMeta
}

type TLBIndexingSwitchMsgBuilder struct {
	sendTime     VTimeInSec
	src, dst     Port
	interleaving uint64
	switchTo     TLBIndexingSwitch
}

func (b TLBIndexingSwitchMsgBuilder) WithSendTime(t VTimeInSec) TLBIndexingSwitchMsgBuilder {
	b.sendTime = t
	return b
}

func (b TLBIndexingSwitchMsgBuilder) WithSrc(src Port) TLBIndexingSwitchMsgBuilder {
	b.src = src
	return b
}

func (b TLBIndexingSwitchMsgBuilder) WithDst(dst Port) TLBIndexingSwitchMsgBuilder {
	b.dst = dst
	return b
}

func (b TLBIndexingSwitchMsgBuilder) WithSwitchTo(
	switchTo TLBIndexingSwitch) TLBIndexingSwitchMsgBuilder {
	b.switchTo = switchTo
	return b
}

func (b TLBIndexingSwitchMsgBuilder) WithInterleaving(
	interleaving uint64) TLBIndexingSwitchMsgBuilder {
	b.interleaving = interleaving
	return b
}

func (b TLBIndexingSwitchMsgBuilder) Build() *TLBIndexingSwitchMsg {
	r := &TLBIndexingSwitchMsg{}
	r.ID = GetIDGenerator().Generate()
	r.Src = b.src
	r.Dst = b.dst
	r.SendTime = b.sendTime
	r.TLBInterleaving = b.interleaving
	r.TLBIndexingSwitch = b.switchTo

	return r
}

// CollectedStatsMsg is sent from TLBs/RTUs/MMUs(?) to provide information
// about different parameters in the prevoius epoch.
type CollectedStatsMsg struct {
	MsgMeta
	Data                                         string
	AverageQueueLength                           float64
	RemoteMemAccessesPerWalk, EnqueuedWalks      float64
	NumHits, NumMisses, NumAccesses, StalledReqs uint64
	IncomingReqs, OutgoingReqs                   uint64
	Appid                                        uint32
}

func (c *CollectedStatsMsg) Meta() *MsgMeta {
	return &c.MsgMeta
}

type CollectedStatsMsgBuilder struct {
	sendTime                                VTimeInSec
	src, dst                                Port
	data                                    string
	averageQueueLength                      float64
	remoteMemAccessesPerWalk, enqueuedWalks float64
	hits, misses, accesses, numStalled      uint64
	incomingReqs, outgoingReqs              uint64
	appid                                   uint32
}

func (b CollectedStatsMsgBuilder) WithSendTime(t VTimeInSec) CollectedStatsMsgBuilder {
	b.sendTime = t
	return b
}

func (b CollectedStatsMsgBuilder) WithSrc(src Port) CollectedStatsMsgBuilder {
	b.src = src
	return b
}

func (b CollectedStatsMsgBuilder) WithDst(dst Port) CollectedStatsMsgBuilder {
	b.dst = dst
	return b
}

func (b CollectedStatsMsgBuilder) WithData(
	data string) CollectedStatsMsgBuilder {
	b.data = data
	return b
}

func (b CollectedStatsMsgBuilder) WithAvgQueueLength(
	avgQueueLength float64) CollectedStatsMsgBuilder {
	b.averageQueueLength = avgQueueLength
	return b
}

func (b CollectedStatsMsgBuilder) WithRemoteMemAccessesPerWalk(
	remoteMemAccesses float64) CollectedStatsMsgBuilder {
	b.remoteMemAccessesPerWalk = remoteMemAccesses
	return b
}

func (b CollectedStatsMsgBuilder) WithNumHits(
	numHits uint64) CollectedStatsMsgBuilder {
	b.hits = numHits
	return b
}

func (b CollectedStatsMsgBuilder) WithNumMisses(
	numMisses uint64) CollectedStatsMsgBuilder {
	b.misses = numMisses
	return b
}

func (b CollectedStatsMsgBuilder) WithNumAccesses(
	numAccesses uint64) CollectedStatsMsgBuilder {
	b.accesses = numAccesses
	return b
}

func (b CollectedStatsMsgBuilder) WithAvgEnqueuedWalks(
	enqueuedWalks float64) CollectedStatsMsgBuilder {
	b.enqueuedWalks = enqueuedWalks
	return b
}

func (b CollectedStatsMsgBuilder) WithIncomingReqs(
	incomingReqs uint64) CollectedStatsMsgBuilder {
	b.incomingReqs = incomingReqs
	return b
}

func (b CollectedStatsMsgBuilder) WithNumStalled(
	numStalled uint64) CollectedStatsMsgBuilder {
	b.numStalled = numStalled
	return b
}

func (b CollectedStatsMsgBuilder) WithOutgoingReqs(
	outgoingReqs uint64) CollectedStatsMsgBuilder {
	b.outgoingReqs = outgoingReqs
	return b
}

func (b CollectedStatsMsgBuilder) Build() *CollectedStatsMsg {
	c := &CollectedStatsMsg{}
	c.ID = GetIDGenerator().Generate()
	c.Src = b.src
	c.Dst = b.dst
	c.SendTime = b.sendTime
	c.Data = b.data
	c.AverageQueueLength = b.averageQueueLength
	c.RemoteMemAccessesPerWalk = b.remoteMemAccessesPerWalk
	c.NumHits = b.hits
	c.NumMisses = b.misses
	c.NumAccesses = b.accesses
	c.EnqueuedWalks = b.enqueuedWalks
	c.IncomingReqs = b.incomingReqs
	c.OutgoingReqs = b.outgoingReqs
	c.StalledReqs = b.numStalled

	return c
}

// SendStatsMsg is to be sent from Command Processor to TLBs/RTUs
// requesting for data collected by TLB/RTUs over the last epoch.
type SendStatsMsg struct {
	MsgMeta
}

func (c *SendStatsMsg) Meta() *MsgMeta {
	return &c.MsgMeta
}

type SendStatsMsgBuilder struct {
	sendTime VTimeInSec
	src, dst Port
}

func (b SendStatsMsgBuilder) WithSendTime(t VTimeInSec) SendStatsMsgBuilder {
	b.sendTime = t
	return b
}

func (b SendStatsMsgBuilder) WithSrc(src Port) SendStatsMsgBuilder {
	b.src = src
	return b
}

func (b SendStatsMsgBuilder) WithDst(dst Port) SendStatsMsgBuilder {
	b.dst = dst
	return b
}

func (b SendStatsMsgBuilder) Build() *SendStatsMsg {
	c := &SendStatsMsg{}
	c.ID = GetIDGenerator().Generate()
	c.Src = b.src
	c.Dst = b.dst
	c.SendTime = b.sendTime
	return c
}

// Trigger Messages are sent from TLBs/RTUs/MMUs to CP on detection on possible
// imbalance conditions locally.
type TriggerMsg struct {
	MsgMeta
}

func (c *TriggerMsg) Meta() *MsgMeta {
	return &c.MsgMeta
}

type TriggerMsgBuilder struct {
	sendTime VTimeInSec
	src, dst Port
}

func (b TriggerMsgBuilder) WithSendTime(t VTimeInSec) TriggerMsgBuilder {
	b.sendTime = t
	return b
}

func (b TriggerMsgBuilder) WithSrc(src Port) TriggerMsgBuilder {
	b.src = src
	return b
}

func (b TriggerMsgBuilder) WithDst(dst Port) TriggerMsgBuilder {
	b.dst = dst
	return b
}

func (b TriggerMsgBuilder) Build() *TriggerMsg {
	c := &TriggerMsg{}
	c.ID = GetIDGenerator().Generate()
	c.Src = b.src
	c.Dst = b.dst
	c.SendTime = b.sendTime
	return c
}

type BalanceRestoredMsg struct {
	MsgMeta
}

func (c *BalanceRestoredMsg) Meta() *MsgMeta {
	return &c.MsgMeta
}

type BalanceRestoredMsgBuilder struct {
	sendTime VTimeInSec
	src, dst Port
}

func (b BalanceRestoredMsgBuilder) WithSendTime(t VTimeInSec) BalanceRestoredMsgBuilder {
	b.sendTime = t
	return b
}

func (b BalanceRestoredMsgBuilder) WithSrc(src Port) BalanceRestoredMsgBuilder {
	b.src = src
	return b
}

func (b BalanceRestoredMsgBuilder) WithDst(dst Port) BalanceRestoredMsgBuilder {
	b.dst = dst
	return b
}

func (b BalanceRestoredMsgBuilder) Build() *BalanceRestoredMsg {
	c := &BalanceRestoredMsg{}
	c.ID = GetIDGenerator().Generate()
	c.Src = b.src
	c.Dst = b.dst
	c.SendTime = b.sendTime
	return c
}
