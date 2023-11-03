package ada

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/savour-labs/wallet-hd-chain/config"
	"testing"
)

func TestAdaClient_GetLatestBlockHeight(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)

	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	height, err := client[0].GetLatestBlockHeight()
	if err != nil {
		panic(err)
	}
	fmt.Println("height:", height)
}

func TestAdaClient_GetAccountBalance(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)

	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	accountBalance, err := client[0].GetAccountBalance("addr1qx46npjjkvwe8rdu0s8w350jk8escgp0we8wkjxxd0mwx74t4xr99vcajwxmclqwargl9v0npssz7ajwadyvv6lkudaq8zwr53")
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(accountBalance)
	fmt.Println("accountBalance:", string(marshal))
}

func TestAdaClient_GetUtxoByAddress(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)
	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	utxo, err := client[0].GetUtxoByAddress("addr1qx46npjjkvwe8rdu0s8w350jk8escgp0we8wkjxxd0mwx74t4xr99vcajwxmclqwargl9v0npssz7ajwadyvv6lkudaq8zwr53")
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(utxo)
	fmt.Println("utxo:", string(marshal))
}

func TestAdaClient_GetTransactionsByAddress(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)

	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	txList, err := client[0].GetTransactionsByAddress("addr1qx46npjjkvwe8rdu0s8w350jk8escgp0we8wkjxxd0mwx74t4xr99vcajwxmclqwargl9v0npssz7ajwadyvv6lkudaq8zwr53",
		nil, nil, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("交易数量：", len(txList))
	marshal, _ := json.Marshal(txList)
	fmt.Println("txList:", string(marshal))
}

func TestAdaClient_GetTransactionsByHash(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)

	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	transaction, err := client[0].
		GetTransactionsByHash("541033e7972ef8e1da70eb4db90a5f4d7c6f441da68316163eb57bf977544809")
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(transaction)
	fmt.Println("transaction:", string(marshal))
}

func TestAdaClient_SendRawTransaction(t *testing.T) {
	var f = flag.String("c", "../../config.yml", "config path")
	flag.Parse()
	conf, err := config.New(*f)
	if err != nil {
		panic(err)
	}
	fmt.Println(conf.Fullnode.Ada.RPCs[0].RPCURL)

	client, err := NewAdaClient(conf)
	if err != nil {
		panic(err)
	}
	// todo need 签好名的交易str
	signedTx := ""
	txHash, err := client[0].SendRawTransaction(signedTx)
	if err != nil {
		panic(err)
	}
	fmt.Println("txHash:", txHash)
}
