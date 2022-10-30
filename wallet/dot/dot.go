package dot

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
)

type Error struct {
	Code    int64
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("code:%d,msg:%s", e.Code, e.Message)
}

// 组装参数
func GetParams(method string, params interface{}) interface{} {
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

type QueryStorageResponse struct {
}

// 根据地址获取账户余额
func QueryStorage(storageKey string, blockHash string) (*QueryStorageResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("state_queryStorage", map[string]interface{}{
		"keys": storageKey,
		"at":   blockHash,
	}))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return nil, nil
}

type GetAccountNonceResponse struct {
	Nonce float64
}

// 获取 Account Nonce
func GetAccountNonce(params []string) (*GetAccountNonceResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("account_nextIndex", params))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return &GetAccountNonceResponse{
		Nonce: jsonObj.Get("result").Float(),
	}, nil
}

type GetBlockResponse struct {
	Number    string `json:"number"`    //最新区块高度
	StateRoot string `json:"stateRoot"` //最新区块Hash
}

// 获取最新 blockNumber 和 blockHash
func GetBlock() (*GetBlockResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("chain_getBlock", nil))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return &GetBlockResponse{
		Number:    jsonObj.Get("result.block.header.number").String(),
		StateRoot: jsonObj.Get("result.block.header.stateRoot").String(),
	}, nil
}

type GetRuntimeVersionResponse struct {
	SpecName           string  `json:"specName"`           //链名称
	SpecVersion        float64 `json:"specVersion"`        //版本
	TransactionVersion float64 `json:"transactionVersion"` //事物版本
}

// 获取specName，specVersion， transactionVersion 等参数
func GetRuntimeVersion() (*GetRuntimeVersionResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("state_getRuntimeVersion", nil))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return &GetRuntimeVersionResponse{
		SpecName:           jsonObj.Get("result.specName").String(),
		SpecVersion:        jsonObj.Get("result.specVersion").Float(),
		TransactionVersion: jsonObj.Get("result.transactionVersion").Float(),
	}, nil
}

type SubmitExtrinsicResponse struct {
}

type SubmitExtrinsicRequest struct {
	Extrinsic map[string]interface{} `json:"extrinsic"`
}

// 广播交易
func SubmitExtrinsic(req *SubmitExtrinsicRequest) (*SubmitExtrinsicResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("author_submitExtrinsic", req.Extrinsic))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return &SubmitExtrinsicResponse{}, nil
}

type GetTxResponse struct {
	SpecName           string  `json:"specName"`           //链名称
	SpecVersion        float64 `json:"specVersion"`        //版本
	TransactionVersion float64 `json:"transactionVersion"` //事物版本
}

// 根据地址获取交易记录
func GetTx() (*GetRuntimeVersionResponse, error) {
	ctx := context.Background()
	respByte, err := DotRequest(ctx, GetParams("state_getRuntimeVersion", []string{"16ZL8yLyXv3V3L3z9ofR1ovFLziyXaN1DPq4yffMAZ9czzBD"}))
	if err != nil {
		return nil, err
	}
	jsonObj := gjson.ParseBytes(respByte)
	if jsonObj.Get("error.code").Int() != 0 {
		return nil, &Error{
			Code:    jsonObj.Get("error.code").Int(),
			Message: jsonObj.Get("error.message").String(),
		}
	}
	return &GetRuntimeVersionResponse{
		SpecName:           jsonObj.Get("result.specName").String(),
		SpecVersion:        jsonObj.Get("result.specVersion").Float(),
		TransactionVersion: jsonObj.Get("result.transactionVersion").Float(),
	}, nil
}
