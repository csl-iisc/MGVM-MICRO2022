package mem

import (
	"gitlab.com/akita/akita"
)

// WriteBuffer is a place where the write can be transferred at a later time.
type WriteBuffer interface {
	Tick(now akita.VTimeInSec) bool
	CanEnqueue() bool
	Enqueue(write *WriteReq)
	Query(read *ReadReq) *WriteReq
	SetWriteCombineGranularity(size uint64)
}

// NewWriteBuffer creates and returns a default write buffer
func NewWriteBuffer(capacity int, port akita.Port) WriteBuffer {
	return &writeBufferImpl{
		capacity:         capacity,
		port:             port,
		writeCombineSize: 64,
	}
}

type writeBufferImpl struct {
	capacity         int
	writeCombineSize uint64
	port             akita.Port
	buf              []*WriteReq
}

func (b *writeBufferImpl) SetWriteCombineGranularity(size uint64) {
	b.writeCombineSize = size
}

func (b *writeBufferImpl) Tick(now akita.VTimeInSec) bool {
	if len(b.buf) == 0 {
		return false
	}

	err := b.port.Send(b.buf[0])
	if err != nil {
		return false
	}

	b.buf = b.buf[1:]
	return true
}

func (b *writeBufferImpl) CanEnqueue() bool {
	if len(b.buf) >= b.capacity {
		return false
	}
	return true
}

func (b *writeBufferImpl) Enqueue(write *WriteReq) {
	if len(b.buf) >= b.capacity {
		panic("Buffer overflow, please use CanEnqueue before Enqueue")
	}

	for _, w := range b.buf {
		if w.PID != write.PID {
			continue
		}

		if w.Address/b.writeCombineSize == write.Address/b.writeCombineSize {
			b.combineWrite(w, write)
			return
		}
	}

	b.buf = append(b.buf, write)
}

func (b *writeBufferImpl) combineWrite(w, write *WriteReq) {
	if w.DirtyMask == nil {
		b.createDirtyMask(w)
	}

	if write.DirtyMask == nil {
		b.createDirtyMask(write)
	}

	b.expandWrite(w)

	offset := write.Address % b.writeCombineSize
	for i := 0; i < len(write.Data); i++ {
		if write.DirtyMask[i] != true {
			continue
		}

		w.Data[offset+uint64(i)] = write.Data[i]
		w.DirtyMask[offset+uint64(i)] = true
	}
}

func (b *writeBufferImpl) createDirtyMask(w *WriteReq) {
	dm := make([]bool, len(w.Data))
	for i := 0; i < len(w.Data); i++ {
		dm[i] = true
	}
	w.DirtyMask = dm
}

func (b *writeBufferImpl) expandWrite(w *WriteReq) {
	if uint64(len(w.Data)) == b.writeCombineSize {
		return
	}

	offset := w.Address % b.writeCombineSize
	newData := make([]byte, b.writeCombineSize)
	newDirtyMask := make([]bool, b.writeCombineSize)
	copy(newData[offset:], w.Data)
	copy(newDirtyMask[offset:], w.DirtyMask)

	w.Data = newData
	w.DirtyMask = newDirtyMask
}

func (b *writeBufferImpl) Query(read *ReadReq) *WriteReq {
	addr := read.Address
	for _, w := range b.buf {
		if w.PID != read.PID {
			continue
		}

		if addr >= w.Address && addr < w.Address+uint64(len(w.Data)) {
			return w
		}
	}
	return nil
}
