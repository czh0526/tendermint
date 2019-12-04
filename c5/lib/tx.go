package lib

import (
	"bytes"
	"time"

	"github.com/tendermint/tendermint/crypto"
)

type Tx struct {
	Payload   Payload
	Signature []byte
	PubKey    crypto.PubKey
	Sequence  int64
}

// 设置 Payload, Sequence 字段
func NewTx(payload Payload) *Tx {
	return &Tx{
		Payload:  payload,
		Sequence: time.Now().Unix(),
	}
}

// 设置 Signature, PubKey 字段
func (tx *Tx) Sign(priv crypto.PrivKey) error {
	data := tx.Payload.GetSignBytes()
	var err error
	tx.Signature, err = priv.Sign(data)
	tx.PubKey = priv.PubKey()
	return err
}

func (tx *Tx) Verify() bool {
	// 验证签名用户是否是 tx.From
	signer := tx.Payload.GetSigner()
	signerFromKey := tx.PubKey.Address()
	if !bytes.Equal(signer, signerFromKey) {
		return false
	}
	// 使用公钥验证签名是否有效
	data := tx.Payload.GetSignBytes()
	sig := tx.Signature
	return tx.PubKey.VerifyBytes(data, sig)
}
