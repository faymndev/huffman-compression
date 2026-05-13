package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode/utf8"
)

const RAW_FILE = "./example.txt"
const COMPRESSED_FILE = "./example.hfm"

func main() {
	CompressFromPath(RAW_FILE, COMPRESSED_FILE)
	DecompressFromPath(COMPRESSED_FILE)
}

func CompressFromPath(path string, destination string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	result, dict, err := CompressStr(string(data))
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.OpenFile(destination, os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
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

const DICT_LENGTH_SIZE = 4
const DICT_ROW_SIZE = 5

func DecompressFromPath(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	dictLength := binary.BigEndian.Uint32(data[:DICT_LENGTH_SIZE])
	dict := make(map[byte]rune)
	for i := range dictLength {
		start := DICT_LENGTH_SIZE + i*DICT_ROW_SIZE
		r := rune(binary.BigEndian.Uint32(data[start : start+4]))
		v := data[start+4]
		dict[v] = r
	}

	sb := strings.Builder{}
	start := DICT_LENGTH_SIZE + dictLength*DICT_ROW_SIZE
	for _, b := range data[start:] {
		letter, ok := dict[b]
		if ok {
			sb.WriteRune(letter)
		} else {
			log.Fatalf("dictionary could not handle %v", b)
		}
	}

	fmt.Println(sb.String())
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

func CompressStr(input string) (result []byte, dict map[rune]byte, err error) {
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
