package circularqueue

import "errors"

const (
	initSize = 32
)

// circularQueue allocate new memory when necessary.
type circularQueue struct {
	buffer []interface{}
	readableIndex int
	writableIndex int
}

// NewCircularQueue creates a circularQueue.
// Usually, you use this function.
func NewCircularQueue() *circularQueue {
	return NewCircularQueueWithSize(initSize)
}

// NewCircularQueueWithSize creates a circularQueue with 
// a size you specified.
func NewCircularQueueWithSize(s int) *circularQueue {
	return &circularQueue{
		buffer: make([]interface{}, s),
		readableIndex: 0,
		writableIndex: 0,
	}
}

// Len returns item count.
func (b *circularQueue) Len() int {
	if b.IsEmpty() {
		return 0
	}

	var length int
	if b.readableIndex < b.writableIndex {
		length = b.writableIndex - b.readableIndex
	} else if b.readableIndex > b.writableIndex {
		length = len(b.buffer)-b.readableIndex + b.writableIndex
	}
	return length
}

// Push pushes a item into this queue.
// Donot worry if this queue is full.
func (b *circularQueue) Push(m interface{}) {
	b.ensureWritableSpace()
	b.buffer[b.writableIndex] = m
	b.hasWritten()
}

func (b *circularQueue) ensureWritableSpace() {
	if b.isFull() {
		b.makeSpace()
	}
}

func (b *circularQueue) makeSpace() {
	buf := make([]interface{}, 1+cap(b.buffer)*2)
	if b.readableIndex < b.writableIndex {
		length := b.writableIndex - b.readableIndex
		copy(buf, b.buffer[b.readableIndex:b.writableIndex])
		b.readableIndex = 0
		b.writableIndex = length
	} else if b.readableIndex > b.writableIndex {
		length := len(b.buffer)-b.readableIndex + b.writableIndex
		copy(buf, b.buffer[b.readableIndex:len(b.buffer)])
		copy(buf[len(b.buffer)-b.readableIndex:], b.buffer[:b.writableIndex])
		b.readableIndex = 0
		b.writableIndex = length
	} else {
		b.readableIndex = 0
		b.writableIndex = 0
	}
	b.buffer = buf
}

func (b *circularQueue) hasWritten() {
	b.writableIndex++
	if b.writableIndex >= len(b.buffer) {
		if b.readableIndex > 0 {
			b.writableIndex = 0
		}
	}
}

// IsEmpty returns true if this queue if empty.
func (b *circularQueue) IsEmpty() bool {
	return b.readableIndex == b.writableIndex
}

func (b *circularQueue) isFull() bool {
	return (b.readableIndex == 0 && b.writableIndex == len(b.buffer)) || 
			b.writableIndex+1 == b.readableIndex
}

func (b *circularQueue) peek() interface{} {
	return b.buffer[b.readableIndex]
}

func (b *circularQueue) retrieve() {
	b.readableIndex++
	if b.writableIndex >= len(b.buffer) {
		b.writableIndex = 0
	}
	if b.readableIndex >= len(b.buffer) {
		b.readableIndex = 0
	}
}

// Pop pops a item.
func (b *circularQueue) Pop() (interface{}, error) {
	if b.IsEmpty() {
		return nil, errors.New("circularQueue is empty")
	}
	m := b.peek()
	b.retrieve()
	return m, nil
}
