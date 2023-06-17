package tezos

import (
	"fmt"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/savour-labs/wallet-hd-chain/wallet/tezos/types"
)

var errTezosHTTPError = errors.New("Tezos chain http error")

type TezosClient interface {
	GetAccountBalance(address string) (*types.AccountBalance, error)
	GetAccountCounter(address string) (*types.AccountCounter, error)
	GetManagerKey(address string) (*types.AccountManagerKey, error)
	SendRawTransaction(rawTx string) (*types.Transaction, error)
}

type Client struct {
	client *gresty.Client
}

func NewTezosClient(url string) (*Client, error) {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errTezosHTTPError)
		}
		return nil
	})
	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetAccountBalance(address string) (balance *types.AccountBalance, err error) {
	var balanceTmp string
	response, err := c.client.R().
		SetResult(&balanceTmp).
		Get("chains/main/blocks/head/context/contracts/" + address + "/balance")
	if err != nil {
		return nil, fmt.Errorf("cannot get account balance: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account balance fail")
	}
	accountBalance := &types.AccountBalance{
		Chain:   "tezos",
		Coin:    "xtz",
		Address: address,
		Balance: balanceTmp,
	}
	return accountBalance, nil
}

func (c *Client) GetAccountCounter(address string) (*types.AccountCounter, error) {
	var counterTemp string
	response, err := c.client.R().
		SetResult(&counterTemp).
		Get("chains/main/blocks/head/context/contracts/" + address + "/counter")
	if err != nil {
		return nil, fmt.Errorf("cannot get account counter: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get account counter fail")
	}
	accountCounter := &types.AccountCounter{
		Chain:   "tezos",
		Coin:    "xtz",
		Address: address,
		Counter: counterTemp,
	}
	return accountCounter, nil
}

func (c *Client) GetManagerKey(address string) (*types.AccountManagerKey, error) {
	var managerkeyTemp string
	response, err := c.client.R().
		SetResult(&managerkeyTemp).
		Get("chains/main/blocks/head/context/contracts/" + address + "/manager_key")
	if err != nil {
		return nil, fmt.Errorf("cannot get manage key: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get manage key fail")
	}
	managerkey := &types.AccountManagerKey{
		Chain:      "tezos",
		Coin:       "xtz",
		Address:    address,
		ManagerKey: managerkeyTemp,
	}
	return managerkey, nil
}

func (c *Client) SendRawTransaction(rawTx string) (*types.Transaction, error) {
	var transactionHash string
	response, err := c.client.R().
		SetBody(map[string]interface{}{"": rawTx}).
		SetResult(&transactionHash).
		Post("injection/operation?async=true&chain=main")
	if err != nil {
		return nil, fmt.Errorf("Can not send raw transaction: %w", err, "code", response.StatusCode())
	}
	tx := &types.Transaction{
		Chain:  "tezos",
		Coin:   "xtz",
		TxHash: transactionHash,
	}
	return tx, nil
}
