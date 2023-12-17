package arweave

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/everFinance/goar"
	"github.com/everFinance/goar/types"

	"github.com/savour-labs/wallet-chain-node/config"
)

type arweaveClient struct {
	client *goar.Client
}

func (a *arweaveClient) GetLatestBlockHeight() (int64, error) {
	info, err := a.client.GetInfo()
	if err != nil {
		return 0, err
	}
	return info.Height, nil
}

func NewArweaveClient(conf *config.Config) ([]*arweaveClient, error) {
	var clients []*arweaveClient
	for _, rpc := range conf.Fullnode.Arweave.RPCs {
		newClient := goar.NewClient(rpc.RPCURL)
		aclient := &arweaveClient{
			client: newClient,
		}
		clients = append(clients, aclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (a *arweaveClient) GetInfo() (*types.NetworkInfo, error) {
	info, err := a.client.GetInfo()
	if err != nil {
		log.Printf("Get info error: %+v\n", err)
		return nil, err
	}
	return info, nil
}

func (a *arweaveClient) GetAccountBalance(address string) (*big.Float, error) {
	amount, err := a.client.GetWalletBalance(address)
	if err != nil {
		log.Printf("Get wallet balance error: %+v\n", err)
		return nil, err
	}
	return amount, nil
}

func (a *arweaveClient) SendRawTransaction(tx string) (string, error) {
	var transaction *types.Transaction
	unmarshalErr := json.Unmarshal([]byte(tx), &transaction)
	if unmarshalErr != nil {
		log.Printf("Unmarshal transaction error: %+v\n", unmarshalErr)
		return "", unmarshalErr
	}
	id := transaction.ID
	status, code, err := a.client.SubmitTransaction(transaction)
	if err != nil {
		log.Printf("Submit transaction error: %+v\n", err)
		return "", err
	}
	if code != http.StatusOK {

		return "", fmt.Errorf("submit transaction http post error: %v", code)
	}
	if status != "OK" {
		return "", fmt.Errorf("submit transaction error: %s", status)
	}
	return id, nil
}
func (a *arweaveClient) GetTransactionById(txHash string) (*types.Transaction, error) {
	transaction, err := a.client.GetTransactionByID(txHash)
	if err != nil {
		log.Printf("Get transaction error: %+v\n", err)
		return nil, err
	}
	return transaction, nil
}

func (a *arweaveClient) GetLastTransactionID(address string) (string, error) {
	transactionID, err := a.client.GetLastTransactionID(address)
	if err != nil {
		log.Printf("Get last transaction id error: %+v\n", err)
		return "", err
	}
	return transactionID, nil
}

func (a *arweaveClient) GetTransactionByTxHash(txHash string) (*TransactionDetail, error) {
	queryGql := "query {\n  transaction(\n    id: \"%s\"\n  ) {\n        id\n        owner {\n          address\n        }\n        recipient\n        quantity {\n          ar\n        }\n        fee {\n          ar\n        }\n        block {\n          id\n          height\n          timestamp\n        }\n      }\n}"
	queryGql = fmt.Sprintf(queryGql, txHash)
	gqlResult, err := a.client.GraphQL(queryGql)
	if err != nil {
		log.Printf("Get transaction list error: %+v\n", err)
		return nil, err
	}
	fmt.Println(string(gqlResult))
	var result TransactionDetail
	jsonErr := json.Unmarshal(gqlResult, &result)
	if jsonErr != nil {
		log.Printf("Unmarshal transaction list error: %+v\n", jsonErr)
		return nil, jsonErr
	}
	return &result, nil
}

func (a *arweaveClient) GetTransactionListByAddress(address, after string, pageSize uint32) (*TransactionList, error) {
	queryGql := "query {\n  transactions(\n    owners: [\"%s\"]\n    first: %d \n after: \"%s\" \n  ) {\n    pageInfo{ hasNextPage}\n    edges {\n      cursor\n      node {\n        id\n        owner {\n          address\n        }\n        recipient\n        quantity {\n          ar\n        }\n        fee {\n          ar\n        }\n        block {\n          id\n   height\n       timestamp\n        }\n      }\n    }\n  }\n}\n"
	queryGql = fmt.Sprintf(queryGql, address, pageSize, after)
	gqlResult, err := a.client.GraphQL(queryGql)
	if err != nil {
		log.Printf("Get transaction list error: %+v\n", err)
		return nil, err
	}
	var result TransactionList
	jsonErr := json.Unmarshal(gqlResult, &result)
	if jsonErr != nil {
		log.Printf("Unmarshal transaction list error: %+v\n", jsonErr)
		return nil, jsonErr
	}
	return &result, nil
}
