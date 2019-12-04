package c5

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/crypto"
)

var (
	SYSTEM_ISSUER = crypto.Address("KING_OF_TOKEN")
)

type TokenApp struct {
	Accounts map[string]int
}

func NewTokenApp() *TokenApp {
	return &TokenApp{
		Accounts: make(map[string]int),
	}
}

func (app *TokenApp) Transfer(from, to crypto.Address, value int) error {
	if app.Accounts[from.String()] < value {
		return errors.New("balance low")
	}
	app.Accounts[from.String()] -= value
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp) Issue(issuer, to crypto.Address, value int) error {
	if !bytes.Equal(issuer, SYSTEM_ISSUER) {
		return errors.New("invalid issuer")
	}
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp) Dump() {
	fmt.Printf("state => %v \n", app.Accounts)
}
