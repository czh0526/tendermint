package lib

import (
	"github.com/tendermint/tendermint/crypto/merkle"
)

func (app *TokenApp) stateToMap() map[string][]byte {
	stateMap := make(map[string][]byte)
	for addr, val := range app.Accounts {
		balance, err := MarshalBinary(val)
		if err != nil {
			panic(err)
		}
		stateMap[addr] = balance
	}
	return stateMap
}

func (app *TokenApp) getRootHash() []byte {
	hashers := app.stateToMap()
	return merkle.SimpleHashFromMap(hashers)
}

func (app *TokenApp) getProofBytes(addr string) ([]byte, *merkle.SimpleProof) {
	stateMap := app.stateToMap()
	root, proofs, _ := merkle.SimpleProofsFromMap(stateMap)
	return root, proofs[addr]
}
