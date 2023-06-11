package eosio

import (
	"context"
	"fmt"
	"strings"

	eos "github.com/eoscanada/eos-go"
	"github.com/eoscanada/eos-go/ecc"
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

func (e *EosClient) GetTransaction(id string) (*eos.TransactionResp, error) {
	res, err := e.client.GetTransaction(context.Background(), id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *EosClient) GetActions(accountName string, pos int64, offset int64) (*eos.ActionsResp, error) {
	infoResp, err := e.client.GetActions(context.Background(), eos.GetActionsRequest{
		AccountName: eos.AccountName(accountName),
		Pos:         eos.Int64(pos),
		Offset:      eos.Int64(offset),
	})
	if err != nil {
		return nil, err
	}
	return infoResp, nil
}

func (e *EosClient) PushTransaction(consumerToken string, rawTx string) (*eos.PushTransactionFullResp, error) {
	infoResp, err := e.client.PushTransaction(context.Background(), &eos.PackedTransaction{
		Signatures: []ecc.Signature{
			ecc.MustNewSignature(consumerToken),
		},
		Compression:       1,
		PackedTransaction: []byte(rawTx),
	})
	if err != nil {
		return nil, err
	}
	return infoResp, nil
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
