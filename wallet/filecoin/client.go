package filecoin

import (
	"context"
	"net/http"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-jsonrpc"
	lotusapi "github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type Client struct {
	closer    jsonrpc.ClientCloser
	apiClient lotusapi.FullNodeStruct
}

func NewClient(addr, authToken string) (*Client, error) { // addr: 127.0.0.1:1234, token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.5eY7Nokrw0eKn6kah21RH2lawg6rbHN1l9mw2s2YxRA
	var api lotusapi.FullNodeStruct
	headers := http.Header{"Authorization": []string{"Bearer " + authToken}}
	closer, err := jsonrpc.NewMergeClient(context.Background(), "ws://"+addr+"/rpc/v0", "Filecoin", []interface{}{&api.Internal, &api.CommonStruct.Internal}, headers)
	if err != nil {
		return nil, err
	}

	return &Client{
		closer:    closer,
		apiClient: api,
	}, nil
}

func (c *Client) Close() error {
	c.closer()
	return nil
}

// GetBalance get balance with coin name and address
func (c *Client) GetBalance(ctx context.Context, addr string) (types.BigInt, error) {
	add, err := address.NewFromString(addr)
	if err != nil {
		return types.EmptyInt, err
	}

	return c.apiClient.WalletBalance(ctx, add)
}

// GetAddress get wallet default address
func (c *Client) GetAddress(ctx context.Context) (address.Address, error) {
	return c.apiClient.WalletDefaultAddress(ctx)
}

// GetWallet get a new wallet
func (c *Client) GetWallet(ctx context.Context) (address.Address, error) {
	return c.apiClient.WalletNew(ctx, types.KTBLS)
}

// GetAddressFromByte convert address from byte to string
func (c *Client) GetAddressFromByte(addr []byte) (string, error) {
	ret, err := address.NewFromBytes(addr)
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}

// VerifyAddress verify if file-coin address
func (c *Client) VerifyAddress(addr string) bool {
	if _, err := address.NewFromString(addr); err != nil {
		return false
	}

	return true
}

// GetNonce MpoolGetNonce
func (c *Client) GetNonce(ctx context.Context, addr string) (uint64, error) {
	add, err := address.NewFromString(addr)
	if err != nil {
		return 0, err
	}

	return c.apiClient.MpoolGetNonce(ctx, add)
}

// Send tx, return cid
func (c *Client) Send(ctx context.Context, fromAddr, toAddr, amount string) (string, error) {
	fromAddress, err := address.NewFromString(fromAddr)
	if err != nil {
		return "", err
	}
	toAddress, err := address.NewFromString(toAddr)
	if err != nil {
		return "", err
	}
	amountBigInt, err := types.BigFromString(amount)
	if err != nil {
		return "", err
	}

	ret, err := c.apiClient.MarketAddBalance(ctx, fromAddress, toAddress, amountBigInt)
	if err != nil {
		return "", err
	}

	return ret.String(), nil
}
