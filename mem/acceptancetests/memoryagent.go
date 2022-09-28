// Package acceptancetests provides utility data structure definitions for
// writing memory system acceptance tests.
package acceptancetests

import (
	"encoding/binary"
	"log"
	"math/rand"
	"reflect"

	"gitlab.com/akita/akita"
	"gitlab.com/akita/mem"
)

// A MemAccessAgent is a Component that can help testing the cache and the the
// memory controllers by generating a large number of read and write requests.
type MemAccessAgent struct {
	*akita.TickingComponent

	LowModule  akita.Port
	MaxAddress uint64

	WriteLeft       int
	ReadLeft        int
	KnownMemValue   map[uint64][]uint32
	PendingReadReq  map[string]*mem.ReadReq
	PendingWriteReq map[string]*mem.WriteReq

	ToMem akita.Port
}

func (a *MemAccessAgent) checkReadResult(
	read *mem.ReadReq,
	dataReady *mem.DataReadyRsp,
) {
	found := false
	var i int
	var value uint32
	result := BytesToUint32(dataReady.Data)
	for i, value = range a.KnownMemValue[read.Address] {
		if value == result {
			found = true
			break
		}
	}

	if found {
		a.KnownMemValue[read.Address] = a.KnownMemValue[read.Address][i:]
	} else {
		log.Panicf("Mismatch when read 0x%X", read.Address)
	}
}

func (a *MemAccessAgent) Tick(now akita.VTimeInSec) bool {
	madeProgress := false

	madeProgress = a.processMsgRsp(now) || madeProgress

	if a.ReadLeft == 0 && a.WriteLeft == 0 {
		return madeProgress
	}

	if a.shouldRead() {
		madeProgress = a.doRead(now) || madeProgress
	} else {
		madeProgress = a.doWrite(now) || madeProgress
	}

	return madeProgress
}

func (a *MemAccessAgent) processMsgRsp(now akita.VTimeInSec) bool {
	msg := a.ToMem.Retrieve(now)
	if msg == nil {
		return false
	}

	switch msg := msg.(type) {
	case *mem.WriteDoneRsp:
		// write := a.PendingWriteReq[msg.RespondTo]
		// log.Printf("%.12f, agent, write complete, 0x%X\n", now, write.Address)
		delete(a.PendingWriteReq, msg.RespondTo)
		return true
	case *mem.DataReadyRsp:
		req := a.PendingReadReq[msg.RespondTo]
		delete(a.PendingReadReq, msg.RespondTo)
		// log.Printf("%.12f, agent, read complete, 0x%X, %v\n",
		// now, req.Address, msg.Data)
		a.checkReadResult(req, msg)
		return true
	default:
		log.Panicf("cannot process message of type %s", reflect.TypeOf(msg))
	}

	return false
}

func (a *MemAccessAgent) shouldRead() bool {
	if len(a.KnownMemValue) == 0 {
		return false
	}

	if a.ReadLeft == 0 {
		return false
	}

	if a.WriteLeft == 0 {
		return true
	}

	dice := rand.Float64()
	return dice > 0.5
}

func (a *MemAccessAgent) doRead(now akita.VTimeInSec) bool {
	address := a.randomReadAddress()

	if a.isAddressInPendingReq(address) {
		return false
	}

	readReq := mem.ReadReqBuilder{}.
		WithSrc(a.ToMem).
		WithDst(a.LowModule).
		WithAddress(address).
		WithByteSize(4).
		WithPID(1).
		Build()
	readReq.SendTime = now
	err := a.ToMem.Send(readReq)
	if err == nil {
		a.PendingReadReq[readReq.ID] = readReq
		a.ReadLeft--
		// log.Printf("%.12f, agent, read, 0x%X\n", now, address)
		return true
	}
	return false
}

func (a *MemAccessAgent) randomReadAddress() uint64 {
	var addr uint64

	for {
		addr = rand.Uint64() % (a.MaxAddress / 4) * 4
		if _, written := a.KnownMemValue[addr]; written {
			return addr
		}
	}
}

func (a *MemAccessAgent) isAddressInPendingReq(addr uint64) bool {
	return a.isAddressInPendingWrite(addr) || a.isAddressInPendingRead(addr)
}

func (a *MemAccessAgent) isAddressInPendingWrite(addr uint64) bool {
	for _, write := range a.PendingWriteReq {
		if write.Address == addr {
			return true
		}
	}
	return false
}

func (a *MemAccessAgent) isAddressInPendingRead(addr uint64) bool {
	for _, read := range a.PendingReadReq {
		if read.Address == addr {
			return true
		}
	}
	return false
}

func Uint32ToBytes(data uint32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, data)
	return bytes
}

func BytesToUint32(data []byte) uint32 {
	a := uint32(0)
	a += uint32(data[0])
	a += uint32(data[1]) << 8
	a += uint32(data[2]) << 16
	a += uint32(data[3]) << 24
	return a
}

func (a *MemAccessAgent) doWrite(now akita.VTimeInSec) bool {
	address := rand.Uint64() % (a.MaxAddress / 4) * 4
	data := rand.Uint32()

	if a.isAddressInPendingReq(address) {
		return false
	}

	writeReq := mem.WriteReqBuilder{}.
		WithSrc(a.ToMem).
		WithDst(a.LowModule).
		WithAddress(address).
		WithPID(1).
		WithData(Uint32ToBytes(data)).
		Build()
	writeReq.SendTime = now

	err := a.ToMem.Send(writeReq)
	if err == nil {
		a.WriteLeft--
		a.addKnownValue(address, data)
		a.PendingWriteReq[writeReq.ID] = writeReq
		// log.Printf("%.12f, agent, write, 0x%X, %v\n",
		// now, address, writeReq.Data)
		return true
	}
	return false
}

func (a *MemAccessAgent) addKnownValue(address uint64, data uint32) {
	valueList, exist := a.KnownMemValue[address]
	if !exist {
		valueList = make([]uint32, 0)
		a.KnownMemValue[address] = valueList
	}
	valueList = append(valueList, data)
	a.KnownMemValue[address] = valueList
}

func NewMemAccessAgent(engine akita.Engine) *MemAccessAgent {
	agent := new(MemAccessAgent)
	agent.TickingComponent = akita.NewTickingComponent("Agent", engine, 1*akita.GHz, agent)
	agent.ToMem = akita.NewLimitNumMsgPort(agent, 1, "Agent.ToMem")

	agent.ReadLeft = 10000
	agent.WriteLeft = 10000
	agent.KnownMemValue = make(map[uint64][]uint32)
	agent.PendingWriteReq = make(map[string]*mem.WriteReq)
	agent.PendingReadReq = make(map[string]*mem.ReadReq)

	return agent
}
