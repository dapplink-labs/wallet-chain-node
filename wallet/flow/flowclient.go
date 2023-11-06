package flow

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/access/http"
	"github.com/savour-labs/wallet-chain-node/config"
	"log"
)

type flowClient struct {
	client *http.Client
}

func NewFlowClient(conf *config.Config) ([]*flowClient, error) {
	var clients []*flowClient
	for _, rpc := range conf.Fullnode.Flow.RPCs {
		client, newClientErr := http.NewClient(rpc.RPCURL)
		if newClientErr != nil {
			continue
		}
		fclient := &flowClient{
			client: client,
		}
		clients = append(clients, fclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (f *flowClient) GetBalance(address string) (int64, error) {
	ctx := context.Background()
	flowAddress := flow.HexToAddress(address)
	account, err := f.client.GetAccount(ctx, flowAddress)
	if err != nil {
		log.Printf("GetBalance  Error: %+v\n", err)
		panic(err)
	}
	return int64(account.Balance), nil
}

func (f *flowClient) GetTxListByAddress(address string) {
	//ctx := context.Background()
	//identifier := flow.HexToID(address)
	//transaction, err := f.client.GetTransaction(ctx, identifier)

}

func (f *flowClient) GetTxDetailByHash(hash string) (*flow.Transaction, error) {
	ctx := context.Background()
	identifier := flow.HexToID(hash)
	transaction, err := f.client.GetTransaction(ctx, identifier)
	if err != nil {
		log.Printf("GetTxDetailByHash  Error: %+v\n", err)
		panic(err)
	}
	return transaction, nil
}

func (f *flowClient) GetGasPrice() {
	//f.client.
}

func (f *flowClient) SendTx(txStr string) error {
	ctx := context.Background()
	var txReq flow.Transaction
	unmarshalErr := json.Unmarshal([]byte(txStr), &txReq)
	if unmarshalErr != nil {
		log.Printf("flow tx unmarshal  Error: %+v\n", unmarshalErr)
		panic(unmarshalErr)
	}
	err := f.client.SendTransaction(ctx, txReq)
	if err != nil {
		log.Printf("SendTransaction  Error: %+v\n", err)
		panic(err)
	}
	return err

}

func (f *flowClient) GetLatestBlockHeight() (int64, error) {
	ctx := context.Background()
	blockHeader, err := f.client.GetLatestBlockHeader(ctx, true)
	if err != nil {
		log.Printf("GetLatestBlockHeight  Error: %+v\n", err)
		panic(err)
	}
	height := blockHeader.Height
	return int64(height), nil
}
