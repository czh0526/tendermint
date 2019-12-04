package main

import (
	"encoding/json"
	"fmt"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)

type CounterApp struct {
	types.BaseApplication
	Value int
}

func NewCounterApp() *CounterApp {
	return &CounterApp{}
}

func (app *CounterApp) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	rsp.Key = []byte("counter")
	rsp.Value = []byte(fmt.Sprintf("%d", app.Value))
	rsp.Log = fmt.Sprintf("value: %d", app.Value)
	return
}

func (app *CounterApp) InitChain(req types.RequestInitChain) (rsp types.ResponseInitChain) {
	var state map[string]int
	err := json.Unmarshal(req.AppStateBytes, &state)
	if err != nil {
		panic(err)
	}
	app.Value = state["counter"]
	return
}

func (app *CounterApp) CheckTx(req types.RequestCheckTx) (rsp types.ResponseCheckTx) {
	if req.GetTx()[0] < 0x04 {
		return
	}

	rsp.Code = 1
	rsp.Log = "bad tx rejected"
	return
}

func (app *CounterApp) DeliverTx(req types.RequestDeliverTx) (rsp types.ResponseDeliverTx) {
	switch req.GetTx()[0] {
	case 0x01:
		app.Value += 1
	case 0x02:
		app.Value -= 1
	case 0x03:
		app.Value = 0
	}
	rsp.Log = fmt.Sprintf("state updated: %d", app.Value)
	return
}

func main() {
	app := NewCounterApp()
	svr, err := server.NewServer("localhost:26658", "socket", app)
	if err != nil {
		panic(err)
	}

	svr.Start()
	defer svr.Stop()
	fmt.Println("abci server started ...")

	select {}
}
