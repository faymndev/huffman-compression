package main

// less than comparator
type Comparator[T any] func(a, b T) bool

type Heap[T any] struct {
	data       []T
	comparator Comparator[T]
}

func NewHeap[T any](comparator Comparator[T]) *Heap[T] {
	var zero T
	heap := &Heap[T]{
		data:       []T{zero},
		comparator: comparator,
	}
	return heap
}

func (h *Heap[T]) Length() int {
	return len(h.data) - 1
}

func (h *Heap[T]) Push(value T) {
	h.data = append(h.data, value)

	// bubble up
	i := len(h.data) - 1
	for i/2 > 0 {
		parent := i / 2
		if h.comparator(h.data[i], h.data[parent]) {
			temp := h.data[i]
			h.data[i] = h.data[parent]
			h.data[parent] = temp
			i = parent
		} else {
			break
		}
	}
}

func (h *Heap[T]) Pop() (result T, ok bool) {
	if len(h.data) == 1 {
		ok = false
		return
	} else if len(h.data) == 2 {
		result = h.data[1]
		h.data = h.data[:len(h.data)-1]
		ok = true
		return
	}

	result = h.data[1]

	// move the last node to the top
	h.data[1] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]
	ok = true

	// bubble down while we have a left node
	i := 1
	for i*2 < len(h.data) {
		if i*2+1 < len(h.data) && h.comparator(h.data[i*2+1], h.data[i]) && h.comparator(h.data[i*2+1], h.data[i*2]) {
			// do we have a right node and is it smaller (than the current and left node)?
			temp := h.data[i]
			h.data[i] = h.data[i*2+1]
			h.data[i*2+1] = temp
			i = i*2 + 1
		} else if h.comparator(h.data[i*2], h.data[i]) {
			// is the left node smaller than the current node?
			temp := h.data[i]
			h.data[i] = h.data[i*2]
			h.data[i*2] = temp
			i = i * 2
		} else {
			// we're done bubbling down
			break
		}
	}

	return
}
