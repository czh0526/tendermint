package main

import (
	"fmt"

	"github.com/czh0526/tendermint/c5"
	"github.com/czh0526/tendermint/c5/lib"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var (
	issuer        = secp256k1.GenPrivKey()
	SYSTEM_ISSUER = issuer.PubKey().Address()
)

func main() {
	app := c5.NewTokenApp()
	p1 := secp256k1.GenPrivKey()
	p2 := secp256k1.GenPrivKey()

	app.Issue(c5.SYSTEM_ISSUER, p1.PubKey().Address(), 1000)
	app.Transfer(p1.PubKey().Address(), p2.PubKey().Address(), 100)
	app.Dump()
	fmt.Println()
	fmt.Println()

	// 测试 Issue Transaction
	var (
		txIssue       *lib.Tx
		rawIssue      []byte
		issueReceived lib.Tx
		err           error
	)
	txIssue = lib.NewTx(lib.NewIssuePayload(
		issuer.PubKey().Address(),
		p1.PubKey().Address(),
		1000,
	))
	txIssue.Sign(issuer)

	rawIssue, err = lib.MarshalBinary((txIssue))
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue tx encoded => %x \n", rawIssue)

	err = lib.UnmarshalBinary(rawIssue, &issueReceived)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue tx decoded => %+v \n", issueReceived)
	fmt.Printf("validated => %t \n", issueReceived.Verify())

	fmt.Println()
	fmt.Println()

	var (
		txTransfer       *lib.Tx
		rawTransfer      []byte
		transferReceived lib.Tx
	)
	txTransfer = lib.NewTx(lib.NewTransferPayload(
		issuer.PubKey().Address(),
		p1.PubKey().Address(),
		1000,
	))
	txTransfer.Sign(issuer)
	rawTransfer, err = lib.MarshalBinary(txTransfer)
	if err != nil {
		panic(err)
	}
	fmt.Printf("transfer tx encoded => %x \n", rawTransfer)

	err = lib.UnmarshalBinary(rawTransfer, &transferReceived)
	if err != nil {
		panic(err)
	}
	fmt.Printf("transfer tx decoded => %+v \n", transferReceived)
	fmt.Printf("validated => %t \n", transferReceived.Verify())
}
