package main

import (
	"fmt"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)

type EzApp struct {
	types.BaseApplication
}

func NewEzApp() *EzApp {
	return &EzApp{}
}

func (app *EzApp) InitChain(req types.RequestInitChain) (rsp types.ResponseInitChain) {
	fmt.Printf("initchain => %+v \n", req)
	return
}

func (app *EzApp) Info(req types.RequestInfo) (rsp types.ResponseInfo) {
	fmt.Printf("info => %+v \n", req)
	return
}

func (app *EzApp) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	fmt.Printf("query => %+v \n", req)
	return
}

func (app *EzApp) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	fmt.Printf("check_tx => %+v \n", req)
	return
}

func (app *EzApp) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	fmt.Printf("deliver_tx => %+v \n", req)
	return
}

func (app *EzApp) BeginBlock(req types.RequestBeginBlock) (rsp types.ResponseBeginBlock) {
	fmt.Printf("begin_block => %+v \n", req)
	return
}

func (app *EzApp) EndBlock(req types.RequestEndBlock) (rsp types.ResponseEndBlock) {
	fmt.Printf("end_block => %+v \n", req)
	return
}

func (app *EzApp) Commit() (rsp types.ResponseCommit) {
	fmt.Printf("commit => \n")
	return
}

func main() {
	app := NewEzApp()
	svr, err := server.NewServer(":26658", "socket", app)
	if err != nil {
		panic(err)
	}

	svr.Start()
	defer svr.Stop()
	fmt.Println("abci server started.")

	select {}
}
