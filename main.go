package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func main() {
	input := "Hello World"

	compressed, root, err := Compress(input)
	if err != nil {
		panic(err)
	}

	originalSize := float64(len(input) * 8)
	compressedSize := float64(len(compressed) * 8)
	savings := 100 * (compressedSize - originalSize) / originalSize

	fmt.Printf("Original   %d bytes\n", len(input))
	fmt.Printf("Compressed %d bytes\n", len(compressed))
	fmt.Printf("%.2f%% reduction in size\n", savings)

	fmt.Printf("Decompressed: %s\n", Decompress(compressed, root))
}

func Decompress(compressed []byte, root *HuffmanNode) string {
	var sb strings.Builder

	currentNode := root
	br := NewBitReader(compressed)

	for br.HasNext() {
		b := br.Next()

		if b == 0 {
			currentNode = currentNode.LeftNode
		} else {
			currentNode = currentNode.RightNode
		}

		if currentNode.IsLeaf() {
			sb.WriteRune(currentNode.Data)
			currentNode = root
		}
	}

	if currentNode.IsLeaf() {
		sb.WriteRune(currentNode.Data)
		currentNode = root
	}

	return sb.String()
}

func Compress(input string) ([]byte, *HuffmanNode, error) {
	heap := NewHeap(func(a, b *HuffmanNode) bool {
		return a.Frequency < b.Frequency
	})

	frequencies := GetFrequencies(input)
	for key, value := range frequencies {
		heap.Push(NewHuffmanNode(key, value, nil, nil))
	}

	for heap.Length() > 1 {
		nodeA, ok := heap.Pop()
		if !ok {
			return nil, nil, errors.New("Expected a smallest node")
		}
		nodeB, ok := heap.Pop()
		if !ok {
			return nil, nil, errors.New("Expected a smallest node")
		}
		heap.Push(NewHuffmanNode(0, 0, nodeA, nodeB))
	}

	root, ok := heap.Pop()
	if !ok {
		return nil, nil, errors.New("Expected a root node")
	}

	runeToCode := make(map[rune]Code)
	Traverse(root, []byte{}, &runeToCode)

	bitWriter := NewBitWriter(10)
	for _, c := range input {
		code := runeToCode[c]
		bitWriter.Write(code.Value, code.Length)
	}
	bitWriter.Flush()

	return bitWriter.Data, root, nil
}

func Traverse(node *HuffmanNode, code []byte, dict *map[rune]Code) {
	if node == nil {
		return
	}

	Traverse(node.LeftNode, append(slices.Clone(code), 0), dict)
	Traverse(node.RightNode, append(slices.Clone(code), 1), dict)

	// leaf node (character)
	if node.LeftNode == nil && node.RightNode == nil {
		// i think real huffman compression algos (that work on ascii anyways) tend to use arrays, such at index 0 = a
		// so array[a] = Code
		(*dict)[node.Data] = NewCode(code)
	}
}

type Code struct {
	Value  byte
	Length int
}

func NewCode(code []byte) Code {
	var value byte = 0
	for _, v := range code {
		value = (value << 1) | v
	}
	return Code{
		Value:  value,
		Length: len(code),
	}
}

func GetFrequencies(raw string) map[rune]int {
	frequencies := make(map[rune]int)
	for _, letter := range raw {
		frequencies[letter]++
	}
	return frequencies
}

type HuffmanNode struct {
	Data      rune
	Frequency int
	LeftNode  *HuffmanNode
	RightNode *HuffmanNode
}

func NewHuffmanNode(data rune, frequency int, leftNode *HuffmanNode, rightNode *HuffmanNode) *HuffmanNode {
	node := &HuffmanNode{
		Data:      data,
		Frequency: frequency,
		LeftNode:  leftNode,
		RightNode: rightNode,
	}

	if leftNode != nil {
		node.Frequency += leftNode.Frequency
	}
	if rightNode != nil {
		node.Frequency += rightNode.Frequency
	}

	return node
}

func (h *HuffmanNode) IsLeaf() bool {
	return h.LeftNode == nil && h.RightNode == nil
}
