package main

type BitWriter struct {
	reservoir byte
	bitCount  int
	Data      []byte
}

func NewBitWriter(sizeHint int) *BitWriter {
	return &BitWriter{
		Data: make([]byte, 0, sizeHint),
	}
}

func (b *BitWriter) Write(code byte, length int) {
	if length <= 0 || length > 8 {
		return
	}

	// clear out garbage bits above the requested length
	// you can't hold more than 8 bits anyways
	// ex) if length is 3, mask is 00001000 - 1 = 00000111
	if length < 8 {
		code &= (1 << length) - 1
	}

	spaceLeft := 8 - b.bitCount
	if length <= spaceLeft {
		// case 1: new bits fit entirely inside of the reservoir
		b.reservoir |= code << (spaceLeft - length)
		b.bitCount += length

		if b.bitCount == 8 {
			b.Data = append(b.Data, b.reservoir)
			b.reservoir = 0
			b.bitCount = 0
		}
	} else {
		// case 2: bits would overflow outside of the reservoir
		overflowBits := b.bitCount + length - 8

		// fill up reservoir and push it
		b.reservoir |= code >> overflowBits
		b.Data = append(b.Data, b.reservoir)

		// set overflow bits
		b.reservoir = code << (8 - overflowBits)
		b.bitCount = overflowBits
	}
}

func (b *BitWriter) Flush() {
	if b.bitCount <= 0 {
		return
	}
	b.Data = append(b.Data, b.reservoir)
	b.reservoir = 0
	b.bitCount = 0
}
