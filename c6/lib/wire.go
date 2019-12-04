package lib

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

var codec = amino.NewCodec()

func init() {
	codec.RegisterInterface((*Payload)(nil), nil)
	codec.RegisterConcrete(&IssuePayload{}, "tx/issue", nil)
	codec.RegisterConcrete(&TransferPayload{}, "tx/transfer", nil)
	codec.RegisterInterface((*crypto.PubKey)(nil), nil)
	codec.RegisterConcrete(&secp256k1.PubKeySecp256k1{}, "secp256k1/pubkey", nil)
	codec.RegisterInterface((*crypto.PrivKey)(nil), nil)
	codec.RegisterConcrete(&secp256k1.PrivKeySecp256k1{}, "secp256k1/privkey", nil)
}

func MarshalJSON(o interface{}) ([]byte, error) {
	return codec.MarshalJSON(o)
}

func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return codec.UnmarshalJSON(bz, ptr)
}

func MarshalBinary(o interface{}) ([]byte, error) {
	return codec.MarshalBinaryBare(o)
}

func UnmarshalBinary(bz []byte, ptr interface{}) error {
	return codec.UnmarshalBinaryBare(bz, ptr)
}
