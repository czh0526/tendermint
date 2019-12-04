package main

import (
	"fmt"

	abci_server "github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)

func main() {
	app := types.NewBaseApplication()
	svr, err := abci_server.NewServer("localhost:26658", "socket", app)
	if err != nil {
		panic(err)
	}

	if err = svr.Start(); err != nil {
		panic(err)
	}
	defer svr.Stop()

	fmt.Println("abci server started ...")

	select {}
}
