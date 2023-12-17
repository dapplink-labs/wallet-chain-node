package arweave

import (
	"flag"
	"fmt"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/savour-labs/wallet-chain-node/config"

	"testing"
)

func TestArweaveClient_GetAccountBalance(t *testing.T) {
	client := getClient()
	balance, err := client.GetAccountBalance("WioN6R1q2s38aZTV8JaG4POYeg1D9qHIjlymya2MWHk")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(balance)
}

func TestArweaveClient_GetInfo(t *testing.T) {
	client := getClient()
	info, err := client.GetInfo()
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(info)
}

func TestArweaveClient_GetLastTransactionID(t *testing.T) {
	client := getClient()
	id, err := client.GetLastTransactionID("du8Kp2JtT9FbtxvYYVVbpHi4lV_fxy7lViz6U4WHV-M")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(id)
}

func TestArweaveClient_GetTransactionById(t *testing.T) {
	client := getClient()
	id, err := client.GetTransactionByTxHash("du8Kp2JtT9FbtxvYYVVbpHi4lV_fxy7lViz6U4WHV-M")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(id)
}

func TestArweaveClient_GetTransactionListByAddress(t *testing.T) {
	client := getClient()
	address, _ := client.GetTransactionListByAddress("Ar-A703DqqY7IhOZLrV1X3JoFvtdfHGhceRlhx2KyRg", "", 10)
	utils.PrettyPrint(address)
}

func getClient() *arweaveClient {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Arweave.RPCs[0].RPCURL)
	clients, err := NewArweaveClient(conf)
	if err != nil {
		panic(err)
	}
	return clients[0]
}
