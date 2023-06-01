package tezos

import (
	"fmt"
	"github.com/SavourDao/savour-hd/wallet/tezos/types"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var errEigenHTTPError = errors.New("Tezos chain http error")

type TezosClient interface {
	getAccountBalance(address string) (*types.AccountBalance, error)
	getAccountCounter(address string) (*types.AccountCounter, error)
	getManagerKey(address string) (*types.AccountManagerKey, error)
}

type Client struct {
	client *gresty.Client
}

func NewTezosClient(url string) *Client {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errEigenHTTPError)
		}
		return nil
	})
	return &Client{
		client: client,
	}
}

func (c *Client) getAccountBalance(address string) (balance *types.AccountBalance, err error) {
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
		Addres:  address,
		Balance: balanceTmp,
	}
	return accountBalance, nil
}

func (c *Client) getAccountCounter(address string) (*types.AccountCounter, error) {
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
		Addres:  address,
		Counter: counterTemp,
	}
	return accountCounter, nil
}

func (c *Client) getManagerKey(address string) (*types.AccountManagerKey, error) {
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
		Addres:     address,
		ManagerKey: managerkeyTemp,
	}
	return managerkey, nil
}
