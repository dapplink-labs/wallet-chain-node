package flows

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	gresty "github.com/go-resty/resty/v2"
	"github.com/onflow/flow-go-sdk"
	flow_http "github.com/onflow/flow-go-sdk/access/http"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/wallet/flows/subgraph"
	"github.com/savour-labs/wallet-chain-node/wallet/flows/subgraph/transaction_detail"
	"github.com/savour-labs/wallet-chain-node/wallet/flows/subgraph/transactions"
)

var errSubGraphHTTPError = errors.New("SubGraph http error")

type flowClient struct {
	client   *flow_http.Client
	sgClient *gresty.Client
}

func NewFlowClient(conf *config.Config) ([]*flowClient, error) {
	var clients []*flowClient
	for _, rpc := range conf.Fullnode.Flow.RPCs {
		client, newClientErr := flow_http.NewClient(rpc.RPCURL)
		if newClientErr != nil {
			continue
		}
		grestyClient := gresty.New()
		grestyClient.SetBaseURL(FlowSubGraphUrl)
		grestyClient.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
			statusCode := r.StatusCode()
			if statusCode >= 400 {
				method := r.Request.Method
				url := r.Request.URL
				return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errSubGraphHTTPError)
			}
			return nil
		})
		fclient := &flowClient{
			client:   client,
			sgClient: grestyClient,
		}
		clients = append(clients, fclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (f *flowClient) GetBalance(address string, proposerKeyIndex uint64) (int64, uint64, error) {
	ctx := context.Background()
	flowAddress := flow.HexToAddress(address)
	account, err := f.client.GetAccountAtLatestBlock(ctx, flowAddress)
	sequenceNumber := account.Keys[proposerKeyIndex].SequenceNumber
	if err != nil {
		log.Printf("GetBalance  Error: %+v\n", err)
		panic(err)
	}
	return int64(account.Balance), sequenceNumber, nil
}

func (f *flowClient) GetTxListByAddress(address string, limit, offset uint64) (*transactions.FlowTransactionsResp, error) {
	var resultTxList transactions.FlowTransactionsResp
	graphRequest := subgraph.GetTransactionsByAddress(address, limit, offset)
	response, err := f.sgClient.R().SetResult(&resultTxList).SetBody(graphRequest).Post("")
	if err != nil {
		log.Printf("GetTxListByAddress  Error: %+v\n", err)
		panic(err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get txList fail")
	}
	return &resultTxList, nil
}

func (f *flowClient) GetTxDetailByHash(hash string) (*transaction_detail.FlowTransactionDetailResp, error) {
	var resultTxDetail transaction_detail.FlowTransactionDetailResp
	graphRequest := subgraph.GetTransactionDetailQuery(hash)
	response, err := f.sgClient.R().SetResult(&resultTxDetail).SetBody(graphRequest).Post("")
	if err != nil {
		log.Printf("GetTxDetailByHash  Error: %+v\n", err)
		panic(err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get tx detail fail")
	}
	return &resultTxDetail, nil
}

func (f *flowClient) GetGasLimit() uint64 {
	return FlowComputeLimit
}

func (f *flowClient) SendTx(txStr string) error {
	ctx := context.Background()
	var txReq flow.Transaction
	unmarshalErr := json.Unmarshal([]byte(txStr), &txReq)
	if unmarshalErr != nil {
		log.Printf("flows tx unmarshal  Error: %+v\n", unmarshalErr)
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
