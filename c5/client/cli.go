package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/czh0526/tendermint/c5/lib"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/crypto"
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
	fmt.Printf("issuer address = %v \n", wallet.GetAddress("issuer"))
}

func query(label string) {
	wallet := lib.LoadWallet("./wallet")
	ret, err := cli.ABCIQuery("", wallet.GetAddress(label))
	if err != nil {
		panic(err)
	}

	var buf []byte
	buf = ret.Response.GetKey()
	addr := crypto.Address(buf)

	buf = ret.Response.GetValue()
	value, err := binary.ReadUvarint(bytes.NewReader(buf))
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v => %d \n", addr, value)
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
