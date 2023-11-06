package flow

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/savour-labs/wallet-chain-node/config"
	"testing"
)

func getClient() *flowClient {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Flow.RPCs[0].RPCURL)
	clients, err := NewFlowClient(conf)
	if err != nil {
		panic(err)
	}
	return clients[0]
}

func TestFlowClient_GetLatestBlockHeight(t *testing.T) {
	client := getClient()
	height, err := client.GetLatestBlockHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println(height)
}

func TestFlowClient_GetTxDetailByHash(t *testing.T) {
	client := getClient()
	txDetail, err := client.GetTxDetailByHash("29d7edc0be31be98898a3457abe001c8d46de7b3b011d8d256f39b78aae5e05e")
	if err != nil {
		panic(err)
	}
	printJsonStr(txDetail)
}

func TestFlowClient_GetBalance(t *testing.T) {
	client := getClient()
	balance, err := client.GetBalance("0x1e3c78c6d580273b")
	if err != nil {
		panic(err)
	}
	printJsonStr(balance)
}

func printJsonStr(param interface{}) {
	marshal, _ := json.Marshal(param)
	fmt.Println(string(marshal))
}
