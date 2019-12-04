package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
)

type CounterApp struct {
	types.BaseApplication
	Value   int
	Version int64
	Hash    string
	History map[int64]int
}

func NewCounterApp() *CounterApp {
	app := &CounterApp{
		History: make(map[int64]int),
	}
	bz, err := ioutil.ReadFile("./counter.state")
	if err != nil {
		return app
	}

	err = json.Unmarshal(bz, app)
	if err != nil {
		return app
	}

	return app
}

func (app *CounterApp) Info(req types.RequestInfo) (rsp types.ResponseInfo) {
	fmt.Printf("info ==> LastBlockHeight = %v, LastBlockAppHash = %v \n", app.Version, app.Hash)
	rsp.LastBlockHeight = app.Version
	rsp.LastBlockAppHash = []byte(app.Hash)
	return
}

func (app *CounterApp) BeginBlock(req types.RequestBeginBlock) (rsp types.ResponseBeginBlock) {
	fmt.Printf("begin_block => %v \n", req.GetHeader().Height)
	return
}

func (app *CounterApp) EndBlock(req types.RequestEndBlock) (rsp types.ResponseEndBlock) {
	fmt.Printf("end_block => %v \n", req.GetHeight())
	return
}

// 每构建一个区块，调用一次 Commit() 函数
func (app *CounterApp) Commit() (rsp types.ResponseCommit) {
	app.Version += 1
	app.History[app.Version] = app.Value
	app.Hash = fmt.Sprintf("%d", app.Value)

	bz, err := json.Marshal(app)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile("./counter.state", bz, 0644)
	rsp.Data = []byte(app.Hash)
	return
}

func (app *CounterApp) Query(req types.RequestQuery) (rsp types.ResponseQuery) {
	ver := req.Height
	if ver == 0 {
		ver = app.Version
	}

	rsp.Key = []byte("counter")
	value := app.History[ver]
	rsp.Value = []byte(fmt.Sprintf("%d", value))
	rsp.Log = fmt.Sprintf("value@%d: %d", ver, value)
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
