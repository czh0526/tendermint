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

	txIssue := lib.NewTx(lib.NewIssuePayload(
		issuer.PubKey().Address(),
		p1.PubKey().Address(),
		1000,
	))
	txIssue.Sign(issuer)
	fmt.Printf("issue tx => %+v \n", txIssue)
	fmt.Printf("validated => %t \n", txIssue.Verify())

	txTransfer := lib.NewTx(lib.NewTransferPayload(
		p1.PubKey().Address(),
		p2.PubKey().Address(),
		100,
	))
	txTransfer.Sign(p1)
	fmt.Printf("transfer tx => %+v \n", txTransfer)
	fmt.Printf("validated => %t \n", txTransfer.Verify())
}
