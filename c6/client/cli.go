package main

import (
	"errors"
	"fmt"

	lib "github.com/czh0526/tendermint/c6/lib"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/rpc/client"
)

var (
	cli = client.NewHTTP("http://localhost:26657", "/websocket")
)

func main() {
	rootCmd := &cobra.Command{
		Use: "cli",
	}

	walletCmd := &cobra.Command{
		Use: "init-wallet",
		Run: func(cmd *cobra.Command, args []string) {
			initWallet()
		},
	}

	issueCmd := &cobra.Command{
		Use: "issue-tx",
		Run: func(cmd *cobra.Command, args []string) {
			issue()
		},
	}

	transferCmd := &cobra.Command{
		Use: "transfer-tx",
		Run: func(cmd *cobra.Command, args []string) {
			transfer()
		},
	}

	queryCmd := &cobra.Command{
		Use: "query",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("query who ?")
			}
			label := args[0]
			query(label)
			return nil
		},
	}

	rootCmd.AddCommand(walletCmd)
	rootCmd.AddCommand(issueCmd)
	rootCmd.AddCommand(transferCmd)
	rootCmd.AddCommand(queryCmd)
	rootCmd.Execute()
}

func initWallet() {
	wallet := lib.NewWallet()
	wallet.GenPrivKey("issuer")
	wallet.GenPrivKey("Cai.Zhihong")
	wallet.GenPrivKey("Wu.Yanhong")
	wallet.Save("./wallet")
}

func issue() {
	wallet := lib.LoadWallet("./wallet")
	tx := lib.NewTx(lib.NewIssuePayload(
		wallet.GetAddress("issuer"),
		wallet.GetAddress("Cai.Zhihong"),
		1000,
	))
	tx.Sign(wallet.GetPrivKey("issuer"))
	bz, err := lib.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}
	ret, err := cli.BroadcastTxCommit(bz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue ret => %+v \n", ret)
}

func transfer() {
	wallet := lib.LoadWallet("./wallet")
	tx := lib.NewTx(lib.NewTransferPayload(
		wallet.GetAddress("Cai.Zhihong"),
		wallet.GetAddress("Wu.Yanhong"),
		100,
	))
	tx.Sign(wallet.GetPrivKey("Cai.Zhihong"))
	bz, err := lib.MarshalBinary(tx)
	if err != nil {
		panic(err)
	}
	ret, err := cli.BroadcastTxCommit(bz)
	if err != nil {
		panic(err)
	}
	fmt.Printf("issue ret => %+v \n", ret)
}

func query(label string) {
	var err error

	// 发送请求
	wallet := lib.LoadWallet("./wallet")
	ret, err := cli.ABCIQueryWithOptions("", wallet.GetAddress(label), client.ABCIQueryOptions{Prove: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v \n", ret)

	var (
		key         crypto.Address
		balance     int
		simpleProof merkle.SimpleProof
	)
	lib.UnmarshalBinary(ret.Response.GetKey(), &key)
	lib.UnmarshalBinary(ret.Response.GetValue(), &balance)

	proof := ret.Response.GetProof()
	ops := proof.GetOps()
	if len(ops) != 2 {
		fmt.Println("response's Proof field is not 2 size.")
		return
	}
	lib.UnmarshalBinary(ops[0].Data, &simpleProof)
	fmt.Printf("return data: %s = %d \n", key, balance)

	rootHash := ops[1].Data
	kvpair := merkle.KVPair{Key: []byte(key.String()), Value: tmhash.Sum(ret.Response.GetValue())}
	fmt.Printf("root hash = %x \n", rootHash)
	fmt.Printf("leafBytes = %x \n", kvpair.Bytes())
	err = simpleProof.Verify(rootHash, kvpair.Bytes())
	if err != nil {
		fmt.Printf("Proof Verify wrong: %v \n", err)
	} else {
		fmt.Printf("Proof Verify success. \n")
	}
}
