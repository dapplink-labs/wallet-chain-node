package ada

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"testing"

	"github.com/savour-labs/wallet-chain-node/config"
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

func TestAdaClient_GetTxFee(t *testing.T) {
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
	txFee, err := client[0].GetTxFee(1000, 569)
	if err != nil {
		panic(err)
	}
	marshal, _ := json.Marshal(txFee)
	fmt.Println("txFee:", string(marshal))
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

	page := 2
	pageSize := 10
	limit := int64(pageSize)
	offset := int64((page - 1) * pageSize)
	txList, err := client[0].GetTransactionsByAddress("addr1qx46npjjkvwe8rdu0s8w350jk8escgp0we8wkjxxd0mwx74t4xr99vcajwxmclqwargl9v0npssz7ajwadyvv6lkudaq8zwr53",
		&limit, &offset)
	if err != nil {
		panic(err)
	}

	fmt.Println("交易数量：", len(txList.Transactions))
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
	input_value_sum := big.NewInt(0)
	output_value_sum := big.NewInt(0)
	for _, operation := range transaction[0].Transaction.Operations {
		toIntValue := stringToInt(operation.Amount.Value)
		if operation.Type == Input {
			input_value_sum = new(big.Int).Add(input_value_sum, toIntValue)
		}
		if operation.Type == Output {
			output_value_sum = output_value_sum.Add(output_value_sum, toIntValue)
		}
	}
	fmt.Println("input_value_sum：", input_value_sum)
	fmt.Println("output_value_sum：", output_value_sum)
	tx_fee := new(big.Int).Sub(new(big.Int).Abs(input_value_sum), new(big.Int).Abs(output_value_sum))
	fmt.Println("交易费用：", tx_fee.String())
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
