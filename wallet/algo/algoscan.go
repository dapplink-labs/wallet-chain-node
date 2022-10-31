package algo

import (
	"fmt"
	"github.com/SavourDao/savour-hd/config"
	"github.com/go-resty/resty/v2"
	"github.com/syndtr/goleveldb/leveldb/errors"
)

var errTssHTTPError = errors.New("tss http error")

type AlgoClient struct {
	client *resty.Client
}

func NewAlgoScanClient(conf *config.Config) *AlgoClient {
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
	return &AlgoClient{
		client: client,
	}
}

func (c *AlgoClient) GetAccount() error {

	return nil
}

func (c *AlgoClient) GettxByAddress() error {
	return nil
}
