package lib

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

var (
	SYSTEM_ISSUER crypto.Address
)

func init() {
	bytes, _ := hex.DecodeString("B56747AA65EA279C35A7948A3E6145D62ABAB3EA")
	SYSTEM_ISSUER = crypto.Address(bytes)
}

type TokenApp struct {
	types.BaseApplication
	Accounts map[string]int
}

func NewTokenApp() *TokenApp {
	return &TokenApp{
		Accounts: make(map[string]int),
	}
}

func (app *TokenApp) Commit() (rsp types.ResponseCommit) {
	rsp.Data = app.getRootHash()
	return
}

func (app *TokenApp) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	addr := crypto.Address(req.Data)
	rsp.Key, _ = MarshalBinary(addr)
	rsp.Value, _ = MarshalBinary(app.Accounts[addr.String()])
	if req.Prove {
		for key, val := range app.Accounts {
			balance, _ := MarshalBinary(val)
			kvPair := merkle.KVPair{Key: []byte(key), Value: tmhash.Sum(balance)}
			fmt.Printf("%s -> %d 's hash => %x \n", key, val, tmhash.Sum(append([]byte{0}, kvPair.Bytes()...)))
		}
		root, proof := app.getProofBytes(addr.String())
		kvpair := merkle.KVPair{Key: []byte(addr.String()), Value: tmhash.Sum(rsp.Value)}
		fmt.Printf("root hash = %x \n", root)
		fmt.Printf("leafBytes = %x \n", kvpair.Bytes())
		if err := proof.Verify(root, kvpair.Bytes()); err != nil {
			fmt.Printf("test verify error: %v \n", err)
		}

		marshaledProof, err := MarshalBinary(proof)
		if err != nil {
			panic(err)
		}

		proofOps := []merkle.ProofOp{
			{Type: "SimpleProof", Data: marshaledProof},
			{Type: "RootHash", Data: root},
		}
		rsp.Proof = &merkle.Proof{Ops: proofOps}
	}
	return
}

func (app *TokenApp) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	tx, err := app.decodeTx(req.GetTx())
	if err != nil {
		rsp.Code = 1
		rsp.Log = "decode error"
		return
	}

	if !tx.Verify() {
		rsp.Code = 2
		rsp.Log = "verify failed"
		return
	}
	return
}

func (app *TokenApp) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	tx, _ := app.decodeTx(req.GetTx())
	switch tx.Payload.GetType() {
	case "issue":
		pld := tx.Payload.(*IssuePayload)
		err := app.issue(pld.Issuer, pld.To, pld.Value)
		if err != nil {
			rsp.Log = err.Error()
		}
		rsp.Info = "issue tx applied"

	case "transfer":
		pld := tx.Payload.(*TransferPayload)
		err := app.transfer(pld.From, pld.To, pld.Value)
		if err != nil {
			rsp.Log = err.Error()
		}
		rsp.Info = "transfer tx applied"
	}
	return
}

func (app *TokenApp) decodeTx(raw []byte) (*Tx, error) {
	var tx Tx
	err := UnmarshalBinary(raw, &tx)
	return &tx, err
}

func (app *TokenApp) transfer(from, to crypto.Address, value int) error {
	if app.Accounts[from.String()] < value {
		return errors.New("balance low")
	}
	app.Accounts[from.String()] -= value
	app.Accounts[to.String()] += value
	return nil
}

func (app *TokenApp) issue(issuer, to crypto.Address, value int) error {
	if !bytes.Equal(issuer, SYSTEM_ISSUER) {
		return errors.New("invalid issuer")
	}
	app.Accounts[to.String()] += value
	return nil
}
