package lib

import (
	"fmt"
	"io/ioutil"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type Wallet struct {
	Keys map[string]crypto.PrivKey
}

func NewWallet() *Wallet {
	return &Wallet{
		Keys: make(map[string]crypto.PrivKey),
	}
}

func (wallet *Wallet) Save(wfn string) {
	bz, err := codec.MarshalJSONIndent(wallet, "", " ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(wfn, bz, 0644)
}

func LoadWallet(wfn string) *Wallet {
	var wallet Wallet
	bz, err := ioutil.ReadFile(wfn)
	if err != nil {
		panic(err)
	}
	err = codec.UnmarshalJSON(bz, &wallet)
	if err != nil {
		panic(err)
	}

	for label, priv := range wallet.Keys {
		fmt.Printf("%s = %x = %x \n", label, priv.PubKey(), priv.PubKey().Address())
	}
	return &wallet
}

func (wallet *Wallet) GenPrivKey(label string) crypto.PrivKey {
	priv := secp256k1.GenPrivKey()
	wallet.Keys[label] = priv
	return priv
}

func (wallet *Wallet) GetPrivKey(label string) crypto.PrivKey {
	return wallet.Keys[label]
}

func (wallet *Wallet) GetPubKey(label string) crypto.PubKey {
	priv := wallet.Keys[label]
	if priv == nil {
		panic("key not found")
	}
	return priv.PubKey()
}

func (wallet *Wallet) GetAddress(label string) crypto.Address {
	priv := wallet.Keys[label]
	if priv == nil {
		panic("key not found")
	}
	return priv.PubKey().Address()
}
