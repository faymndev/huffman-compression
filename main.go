package main

import (
	"encoding/binary"
	"errors"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

func main() {
	data, err := os.ReadFile("./example.txt")
	if err != nil {
		log.Fatal(err)
	}

	result, dict, err := Compress(string(data))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile("example.hfm", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// dict header:
	// - number of rows
	// - each row has the rune, size of the key, and the key (bytes)
	buf := make([]byte, 4)

	binary.BigEndian.PutUint32(buf, uint32(len(dict)))
	f.Write(buf)

	for key, value := range dict {
		binary.BigEndian.PutUint32(buf, uint32(key))
		f.Write(buf)
		f.Write([]byte{value})
	}
	f.Write(result)
	f.Sync()
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

func Compress(input string) (result []byte, dict map[rune]byte, err error) {
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
	if !ok {
		err = errors.New("Expected a root node")
		return
	}

	dict = make(map[rune]byte)
	Traverse(root, 0b0, &dict)

	result = make([]byte, utf8.RuneCountInString(input))
	for i, r := range input {
		result[i] = dict[r]
	}

	return
}

func Traverse(node *HuffmanNode, code byte, dict *map[rune]byte) {
	if node == nil {
		return
	}

	if node.Data != 0 {
		(*dict)[node.Data] = code
	}

	Traverse(node.LeftNode, code<<1, dict)
	Traverse(node.RightNode, code<<1+0b1, dict)
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
