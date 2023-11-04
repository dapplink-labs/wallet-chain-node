package algo

import (
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/syndtr/goleveldb/leveldb/errors"

	"github.com/savour-labs/wallet-chain-node/config"
)

var errTssHTTPError = errors.New("tss http error")

type AlgoScanClient struct {
	client *resty.Client
}

func NewAlgoScanClient(conf *config.Config) *AlgoScanClient {
	client := resty.New()
	client.SetHostURL(conf.Fullnode.Algo.TpApiUrl)
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errTssHTTPError)
		}
		return nil
	})
	return &AlgoScanClient{
		client: client,
	}
}

func (asc *AlgoScanClient) GetAccount() error {

	return nil
}

func (asc *AlgoScanClient) GettxByAddress() error {
	return nil
}
