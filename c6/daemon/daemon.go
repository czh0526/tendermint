package main

import (
	"fmt"

	"github.com/czh0526/tendermint/c6/lib"
	"github.com/tendermint/tendermint/abci/server"
)

func main() {
	app := lib.NewTokenApp()
	svr, err := server.NewServer(":26658", "socket", app)
	if err != nil {
		panic(err)
	}
	svr.Start()
	defer svr.Stop()

	fmt.Println("token server started.")
	select {}
}
