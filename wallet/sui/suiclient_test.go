package sui

import (
	"flag"
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/utils"

	"github.com/savour-labs/wallet-chain-node/config"
)

func TestSuiClient_GetAllAccountBalance(t *testing.T) {
	client := getClient()

	balance, err := client.GetAllAccountBalance("0x00878369f475a454939af7b84cdd981515b1329f159a1aeb9bf0f8899e00083a")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(balance)
}

func TestSuiClient_GetAccountBalance(t *testing.T) {
	client := getClient()
	balance, err := client.GetAccountBalance("0x00878369f475a454939af7b84cdd981515b1329f159a1aeb9bf0f8899e00083a", "")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(balance)
}

func TestSuiClient_GetTxListByAddress(t *testing.T) {
	client := getClient()
	//"0x7d20dcdb2bca4f508ea9613994683eb4e76e9c4ed371169677c1be02aaf0b58e"
	txList, err := client.GetTxListByAddress("0x95f1baf8c250c06fc2558f2ca5b35b371977f7182d381cf29b0f36f2f9da434a", "YxjRfteuVNyPfJdTf3gZD6grHjUrkTgi8pQKQZqGHyz", 5)
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(txList)
}

func TestSuiClient_GetTxDetailByDigest(t *testing.T) {
	client := getClient()

	txDetail, err := client.GetTxDetailByDigest("Tgc2M6cBMGoYidew1gC2LYwfqQzEBpK2jSAwhCWRCtJ")
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(txDetail)
}

func TestSuiClient_GetGasPrice(t *testing.T) {
	client := getClient()

	gasPrice, err := client.GetGasPrice()
	if err != nil {
		panic(err)
	}
	utils.PrettyPrint(gasPrice)
}

func getClient() *suiClient {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Sui.RPCs[0].RPCURL)
	clients, err := NewSuiClient(conf)
	if err != nil {
		panic(err)
	}
	return clients[0]
}
