package util

import "log"

// A Buffer is a fifo queue for anything
type Buffer interface {
	CanPush() bool
	Push(e interface{})
	Pop() interface{}
	Peek() interface{}
	Capacity() int
	Size() int

	// Remove all elements in the buffer
	Clear()
}

// NewBuffer creates a default buffer object.
func NewBuffer(capacity int) Buffer {
	return &bufferImpl{capacity: capacity}
}

type bufferImpl struct {
	capacity int
	elements []interface{}
}

func (b *bufferImpl) CanPush() bool {
	return len(b.elements) < b.capacity
}

func (b *bufferImpl) Push(e interface{}) {
	if len(b.elements) >= b.capacity {
		log.Panic("buffer overflow")
	}
	b.elements = append(b.elements, e)
}

func (b *bufferImpl) Pop() interface{} {
	if len(b.elements) == 0 {
		return nil
	}

	e := b.elements[0]
	b.elements = b.elements[1:]
	return e
}

func (b *bufferImpl) Peek() interface{} {
	if len(b.elements) == 0 {
		return nil
	}

	return b.elements[0]
}

func (b *bufferImpl) Capacity() int {
	return b.capacity
}

func (b *bufferImpl) Size() int {
	return len(b.elements)
}

func (b *bufferImpl) Clear() {
	b.elements = nil
}
