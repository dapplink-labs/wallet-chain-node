package tezos

import (
	"fmt"
	"github.com/SavourDao/savour-hd/wallet/tezos/types"
	gresty "github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

var errTezosScanHTTPError = errors.New("Tezos chain http error")

type TezosScanClient interface {
	GetTransactionListByAddress(address string) (txList *types.TransactionList, err error)
	// GetTransactionByHash(hash string) (txList *types.TransactionTxHashList, err error)
}

type ScanClient struct {
	client *gresty.Client
}

func NewTezosScanClient(url string) *ScanClient {
	client := gresty.New()
	client.SetHostURL(url)
	client.OnAfterResponse(func(c *gresty.Client, r *gresty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			url := r.Request.URL
			return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errTezosScanHTTPError)
		}
		return nil
	})
	return &ScanClient{
		client: client,
	}
}

//func (sc *ScanClient) GetTransactionListByAddress(address string) (txList *types.TransactionList, err error) {
//	var transactionListTemp *types.TransactionList
//	response, err := sc.client.R().
//		SetResult(transactionListTemp).
//		Get("/explorer/op/" + hash + "")
//	if err != nil {
//		return nil, fmt.Errorf("cannot get transaction detail: %w", err)
//	}
//	if response.StatusCode() != 200 {
//		return nil, errors.New("get transaction detail fail")
//	}
//	return transactionListTemp, nil
//}

func (sc *ScanClient) GetTransactionByHash(hash string) (txList *[]types.TransactionTxHash, err error) {
	var transactionListTemp *[]types.TransactionTxHash
	response, err := sc.client.R().
		SetResult(transactionListTemp).
		Get("/explorer/op/" + hash + "")
	if err != nil {
		return nil, fmt.Errorf("cannot get transaction detail: %w", err)
	}
	if response.StatusCode() != 200 {
		return nil, errors.New("get transaction detail fail")
	}
	return transactionListTemp, nil
}
