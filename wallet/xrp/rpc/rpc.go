package nearrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/SavourDao/savour-hd/wallet/near/types"
	"io/ioutil"
	"net/http"
)

type RRCClient interface {
	DoRpcRequest(method string, params any, result interface{}) error
}

type RpcClient struct {
	URL string
}

func (r RpcClient) RpcRequest(method string, params any, result interface{}) error {
	reqParams := types.RpcRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	jsonStr, _ := json.Marshal(reqParams)
	req, e := http.NewRequest("POST", r.URL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	e = json.Unmarshal(body, &result)
	if e != nil {
		return e
	}
	return nil
}

func (r RpcClient) Request(params string, result interface{}) error {
	fmt.Println(params)
	fmt.Println(r.URL)
	req, e := http.NewRequest("POST", r.URL, bytes.NewBuffer([]byte(params)))
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return e
	}
	e = json.Unmarshal(body, &result)
	if e != nil {
		return e
	}
	return nil
}
