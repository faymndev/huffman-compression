package main

import (
	"slices"
	"testing"
)

func TestHeap(t *testing.T) {
	comparator := func(a, b int) bool {
		return a < b
	}

	cases := [][]int{
		{14, 90, 41, 40, 31, 20},
		{2, 1, 4, 3},
	}

	for _, c := range cases {

		heap := NewHeap(comparator)

		for _, value := range c {
			heap.Push(value)
		}

		slices.Sort(c)
		for i := range c {
			expected := c[i]
			m, ok := heap.Pop()
			if !(ok && m == expected) {
				t.Errorf("expected %d, got %d", expected, m)
			}
		}
	}
}
