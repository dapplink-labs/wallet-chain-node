package dot

import (
	"fmt"
	"github.com/SavourDao/savour-hd/config"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

var errTssHTTPError = errors.New("tss http error")

type dotClient struct {
	client *resty.Client
}

func NewDotClient(conf *config.Config) ([]*dotClient, error) {
	var clients []*dotClient
	for _, rpc := range conf.Fullnode.Eth.RPCs {
		client := resty.New()
		client.SetHostURL(rpc.RPCURL)
		client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
			statusCode := r.StatusCode()
			if statusCode >= 400 {
				method := r.Request.Method
				url := r.Request.URL
				return fmt.Errorf("%d cannot %s %s: %w", statusCode, method, url, errTssHTTPError)
			}
			return nil
		})
		dclient := &dotClient{
			client: client,
		}
		clients = append(clients, dclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (c *dotClient) GetParams(method string, params interface{}) interface{} {
	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"id":      1,
	}
	if params != nil {
		req["params"] = params
	}
	return req
}

func (c *dotClient) GetLatestBlockHeight() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (c *dotClient) GetAccountNonce(params interface{}) (uint64, error) {
	response, err := c.client.R().
		SetBody(c.GetParams("account_nextIndex", params)).
		Post("")
	if err != nil {
		return 0, fmt.Errorf("cannot get signature: %w", err)
	}
	response.Body()
	jsonObj := gjson.ParseBytes(response.Body())
	if jsonObj.Get("error.code").Int() != 0 {
		return 0, nil
	}
	return jsonObj.Get("result").Uint(), nil
}

func (c *dotClient) GetRuntimeVersion() (*RuntimeVersionResponse, error) {
	response, err := c.client.R().
		SetBody(c.GetParams("state_getRuntimeVersion", nil)).
		Post("")
	if err != nil {
		return nil, fmt.Errorf("cannot get result: %w", err)
	}
	jsonObj := gjson.ParseBytes(response.Body())
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, errors.New("GetRuntimeVersion: parse result fail")
	}
	return &RuntimeVersionResponse{
		SpecName:           jsonObj.Get("result.specName").String(),
		SpecVersion:        jsonObj.Get("result.specVersion").Float(),
		TransactionVersion: jsonObj.Get("result.transactionVersion").Float(),
	}, nil
}

func (c *dotClient) SubmitExtrinsic(rawTx string) (string, error) {
	response, err := c.client.R().
		SetBody(c.GetParams("author_submitExtrinsic", rawTx)).
		Post("")
	if err != nil {
		return "", fmt.Errorf("cannot get result: %w", err)
	}
	jsonObj := gjson.ParseBytes(response.Body())
	if jsonObj.Get("error.code").Int() != 0 {
		return "", errors.New("SubmitExtrinsic: parse result fail")
	}
	return jsonObj.Get("hash").Str, nil
}

func (c *dotClient) GetTx(hash string) (*Transaction, error) {
	response, err := c.client.R().
		SetBody(c.GetParams("state_getRuntimeVersion", hash)).
		Post("")
	if err != nil {
		return nil, fmt.Errorf("cannot get result: %w", err)
	}
	jsonObj := gjson.ParseBytes(response.Body())
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, errors.New("GetTx: parse result fail")
	}
	return &Transaction{
		SpecName:           jsonObj.Get("result.specName").String(),
		SpecVersion:        jsonObj.Get("result.specVersion").Float(),
		TransactionVersion: jsonObj.Get("result.transactionVersion").Float(),
	}, nil
}
