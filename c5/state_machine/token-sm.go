package main

import (
	"github.com/czh0526/tendermint/c5"
	"github.com/tendermint/tendermint/crypto"
)

func main() {
	app := c5.NewTokenApp()
	a1 := crypto.Address("TEST_ADDR1")
	a2 := crypto.Address("TEST_ADDR2")

	app.Issue(c5.SYSTEM_ISSUER, a1, 1000)
	app.Transfer(a1, a2, 100)
	app.Dump()
}
