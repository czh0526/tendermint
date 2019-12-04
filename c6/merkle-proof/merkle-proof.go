package main

import (
	"fmt"

	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func sliceDemo() {
	data := [][]byte{
		[]byte("one"), []byte("two"), []byte("three"), []byte("four"),
	}
	for _, d := range data {
		fmt.Printf("%s's hash => %x \n", d, tmhash.Sum(append([]byte{0}, d...)))
	}
	// 生成 Root Hash 和 Proof
	root, proofs := merkle.SimpleProofsFromByteSlices(data)

	// 正确性验证
	err := proofs[0].Verify(root, data[0])
	if err != nil {
		fmt.Printf("data[0] is invalid => err = %v \n", err)
	} else {
		fmt.Println("data[0] is valid")
	}
}

func mapDemo() {
	data := map[string][]byte{
		"tom":   []byte("actor"),
		"mary":  []byte("teacher"),
		"linda": []byte("scientist"),
		"luke":  []byte("fisher"),
	}
	for key, val := range data {
		kvPair := merkle.KVPair{Key: []byte(key), Value: tmhash.Sum([]byte(val))}
		fmt.Printf("%s -> %s 's hash => %x \n", key, val, tmhash.Sum(append([]byte{0}, kvPair.Bytes()...)))
	}

	// 生成 root hash 和 Proof
	root, proofs, keys := merkle.SimpleProofsFromMap(data)
	fmt.Printf("root hash => %x \n", root)
	fmt.Printf("proof for tom => %+v \n", proofs["tom"])
	fmt.Printf("keys sorted => %v \n", keys)

	// 验证
	kvpair := merkle.KVPair{Key: []byte("tom"), Value: tmhash.Sum([]byte("actor"))}
	err := proofs["tom"].Verify(root, kvpair.Bytes())
	if err != nil {
		fmt.Printf("data[0] is invalid => err = %v \n", err)
	} else {
		fmt.Println("data[0] is valid")
	}
}

func main() {
	sliceDemo()
	mapDemo()
}
