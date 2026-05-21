package main

type BitReader struct {
	currentByte    int
	currentByteBit int // tracked from 0 to 7
	data           []byte
	bitsRead       int
}

func NewBitReader(data []byte) *BitReader {
	return &BitReader{
		data: data,
	}
}

func (b *BitReader) GetBitsRead() int {
	return b.bitsRead
}

func (b *BitReader) Next() byte {
	if !b.HasNext() {
		return 0
	}

	data := b.data[b.currentByte]
	b.bitsRead += 1

	// shift to the right such that the current bit is the right most bit now
	// then mask out everything to the left
	var next byte = (data >> (7 - b.currentByteBit)) & 0x01

	// advance the bit pointer
	b.currentByteBit++
	if b.currentByteBit >= 8 {
		b.currentByte += 1
		b.currentByteBit = 0
	}

	return next
}

func (b *BitReader) HasNext() bool {
	return b.currentByte < len(b.data)
}
