package bitcoin

import (
	"fmt"

	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"

	"github.com/savour-labs/wallet-chain-node/wallet/bitcoin/types"
)

/*
 * api docs:
 *	https://www.blockchain.com/explorer/api/blockchain_api
 */
var errBlockChainHTTPError = errors.New("blockchain http error")

type BcClient struct {
	client *gresty.Client
}

func NewBlockChainClient(url string) (*BcClient, error) {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errBlockChainHTTPError)
		}
		return nil
	})
	return &BcClient{
		client: client,
	}, nil
}

func (c *BcClient) GetAccountBalance(address string) (string, error) {
	var accountBalance map[string]*types.AccountBalance
	response, err := c.client.R().
		SetResult(&accountBalance).
		Get("/balance?active=" + address)
	if err != nil {
		return "", fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return "", errors.New("get account balance fail")
	}
	return accountBalance[address].FinalBalance.String(), nil
}

func (c *BcClient) GetAccountUtxo(address string) ([]types.UnspentOutput, error) {
	var utxoUnspentList types.UnspentOutputList
	response, err := c.client.R().
		SetResult(&utxoUnspentList).
		Get("/unspent?active=" + address)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return utxoUnspentList.UnspentOutputs, nil
}

func (c *BcClient) GetTransactionsByAddress(address, pageSize, page string) (*types.Transaction, error) {
	var transactionList types.Transaction
	response, err := c.client.R().
		SetResult(&transactionList).
		Get("/rawaddr/" + address + "?limit=" + pageSize + "&offset=" + page)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transactionList, nil
}

func (c *BcClient) GetTransactionsByHash(txHash string) (*types.TxsItem, error) {
	var transaction types.TxsItem
	response, err := c.client.R().
		SetResult(&transaction).
		Get("/rawtx/" + txHash)
	if err != nil {
		return nil, fmt.Errorf("cannot utxo fail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	return &transaction, nil
}
