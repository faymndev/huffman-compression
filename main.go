package main

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

func main() {
	compressed, dict, err := Compress("Hello World")
	if err != nil {
		panic(err)
	}

	decompressed, err := DecompressOriginal(compressed, dict)
	fmt.Println(decompressed)
}

func DecompressOriginal(input []byte, dict map[rune]byte) (result string, err error) {
	// invert dictionary to go from byte to rune
	inverted := make(map[byte]rune)
	for key, value := range dict {
		inverted[value] = key
	}

	r := strings.Builder{}
	for _, i := range input {
		r.WriteRune(inverted[i])
	}

	result = r.String()
	return
}

func Compress(input string) (result []byte, huffmanNode *HuffmanNode, err error) {
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
			err = errors.New("Expected a smallest node")
			return
		}
		nodeB, ok := heap.Pop()
		if !ok {
			err = errors.New("Expected a smallest node")
			return
		}
		heap.Push(NewHuffmanNode(0, 0, nodeA, nodeB))
	}

	root, ok := heap.Pop()
	huffmanNode = root
	if !ok {
		err = errors.New("Expected a root node")
		return
	}

	dict := make(map[rune]Code)
	Traverse(root, []byte{}, &dict)

	// for _, r := range input {
	// 	// result = append(result, (d[r].Bits))
	// }

	return
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
