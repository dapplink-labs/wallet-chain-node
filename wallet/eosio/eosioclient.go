package eosio

import (
	"context"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
)

type EosClient struct {
	client *eos.API
}

func newClient(url string) (*EosClient, error) {
	return &EosClient{
		client: eos.New(url),
	}, nil
}

func (e *EosClient) ABIBinToJSON(id string) {
	_, err := e.client.ABIBinToJSON(
		context.Background(),
		eos.AccountName(""),
		eos.Name(""),
		eos.HexBytes{},
	)
	if err != nil {
		return
	}
}

func (e *EosClient) ABIJSONToBin(id string) {
	fmt.Println("ABIJSONToBin start")
	infoResp, err := e.client.ABIJSONToBin(
		context.Background(),
		eos.AccountName(""),
		eos.Name(""),
		eos.M{},
	)
	if err != nil {
		fmt.Println("ABIJSONToBin error")
		return
	}
	fmt.Println("ABIJSONToBin res", infoResp)
}

func (e *EosClient) GetTransaction(id string) {
	fmt.Println("GetTransaction start")
	infoResp, err := e.client.GetTransaction(context.Background(), id)
	if err != nil {
		fmt.Println("GetTransaction error")
		return
	}
	fmt.Println("GetTransaction res", infoResp)
}

func (e *EosClient) GetActions() {
	fmt.Println("GetActions start")
	infoResp, err := e.client.GetActions(context.Background(), eos.GetActionsRequest{})
	if err != nil {
		fmt.Println("GetActions error")
		return
	}
	fmt.Println("GetActions res", infoResp)
}

func (e *EosClient) PushTransaction() {
	fmt.Println("PushTransaction start")
	infoResp, err := e.client.PushTransaction(context.Background(), &eos.PackedTransaction{})
	if err != nil {
		fmt.Println("PushTransaction error")
		return
	}
	fmt.Println("PushTransaction res", infoResp)
}

func (e *EosClient) GetAccount(accountName string) (string, error) {
	infoResp, err := e.client.GetAccount(context.Background(), eos.AccountName(accountName))
	if err != nil {
		return "", err
	}
	balanceWithSymbol := infoResp.CoreLiquidBalance.String()
	balanceWithSymbolList := strings.Split(balanceWithSymbol, " ")
	return balanceWithSymbolList[0], nil
}

func (e *EosClient) GetInfo() {
	fmt.Println("GetInfo start")
	infoResp, err := e.client.GetInfo(context.Background())
	if err != nil {
		fmt.Println("GetInfo error")
		return
	}
	fmt.Println("GetInfo res", infoResp)
}
