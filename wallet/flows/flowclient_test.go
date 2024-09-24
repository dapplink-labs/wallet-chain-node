package flows

import (
	"encoding/json"
	"flag"
	"fmt"
	"testing"

	"github.com/savour-labs/wallet-chain-node/config"
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
	txDetail, err := client.GetTxDetailByHash("736466a1261e3caddbd6fc0e143aad0c543235733a58740273a5dac5bfefc031")
	if err != nil {
		panic(err)
	}
	printJsonStr(txDetail)
}

func TestFlowClient_GetBalance(t *testing.T) {
	client := getClient()
	balance, sequenceNumber, err := client.GetBalance("0x1e3c78c6d580273b", 0)
	if err != nil {
		panic(err)
	}
	printJsonStr(balance)
	printJsonStr(sequenceNumber)
}

func printJsonStr(param interface{}) {
	marshal, _ := json.Marshal(param)
	fmt.Println(string(marshal))
}

func TestFlowClient_GetTxListByAddress(t *testing.T) {
	//address := subgraph.GetTransactionsByAddress("0x18eb4ee6b3c026d2", 25, 0)
	//marshal, _ := json.Marshal(address)
	//fmt.Println(string(marshal))
	client := getClient()
	address, err := client.GetTxListByAddress("0x18eb4ee6b3c026d2", 25, 0)
	if err != nil {
		panic(err)
	}
	//utils.PrettyPrint(address)
	printJsonStr(address)
}
