package sui

import (
	"flag"
	"fmt"
	"github.com/block-vision/sui-go-sdk/utils"
	"github.com/savour-labs/wallet-chain-node/config"
	"testing"
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

	client.GetTxListByAddress("0x00878369f475a454939af7b84cdd981515b1329f159a1aeb9bf0f8899e00083a")

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
