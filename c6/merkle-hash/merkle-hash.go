package main

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/merkle"
)

func sliceDemo() {
	var (
		hash []byte
		data [][]byte
	)
	data = [][]byte{
		[]byte("one"), []byte("two"), []byte("three"), []byte("four"),
	}
	hash = merkle.SimpleHashFromByteSlices(data)
	fmt.Printf("[Slice Demo]:\t root hash => %x \n", hash)

	// 改变数据集的次序
	data = [][]byte{
		[]byte("one"), []byte("three"), []byte("four"), []byte("two"),
	}
	hash = merkle.SimpleHashFromByteSlices(data)
	fmt.Printf("[Slice Demo]:\t root hash => %x \n", hash)
}

func mapDemo() {
	var (
		hash []byte
		data map[string][]byte
	)
	data = map[string][]byte{
		"tom":   []byte("actor"),
		"mary":  []byte("teacher"),
		"linda": []byte("scientist"),
		"luke":  []byte("fisher"),
	}
	hash = merkle.SimpleHashFromMap(data)
	fmt.Printf("[Map Demo]:\t root hash => %x \n", hash)

	// 改变数据集的次序
	data = map[string][]byte{
		"mary":  []byte("teacher"),
		"tom":   []byte("actor"),
		"linda": []byte("scientist"),
		"luke":  []byte("fisher"),
	}
	hash = merkle.SimpleHashFromMap(data)
	fmt.Printf("[Map Demo]:\t root hash => %x \n", hash)
}

func main() {
	sliceDemo()
	mapDemo()
}
