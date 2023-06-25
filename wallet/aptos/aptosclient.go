package aptos

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"emperror.dev/errors"
)

type Client struct {
	endpoint   string
	httpClient *http.Client
}

func NewClient(endpoint string) (*Client, error) {
	return &Client{
		endpoint:   endpoint,
		httpClient: new(http.Client),
	}, nil
}

// GetSeq get sequence number from address
func (c *Client) GetSeq(addr string) (int, error) {
	route := fmt.Sprintf("/v1/accounts/%s", addr)

	ret, err := c.doGet(route)
	if err != nil {
		return 0, err
	}

	res := struct {
		SequenceNumber    string `json:"sequence_number"`
		AuthenticationKey string `json:"authentication_key"`
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return 0, err
	}

	return strconv.Atoi(res.SequenceNumber)
}

type EstimateGasPrice struct {
	DeprioritizedGasEstimate int `json:"deprioritized_gas_estimate"`
	GasEstimate              int `json:"gas_estimate"`
	PrioritizedGasEstimate   int `json:"prioritized_gas_estimate"`
}

// GetGasPrice Get estimate gas price
func (c *Client) GetGasPrice() (*EstimateGasPrice, error) {
	route := fmt.Sprintf("/v1/estimate_gas_price")

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := new(EstimateGasPrice)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

type BroadcastTxCont struct {
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Signature               struct {
		Type      string `json:"type"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}
}

// BroadcastTx send tx to the network
func (c *Client) BroadcastTx(tx *BroadcastTxCont) (string, error) {
	route := "/v1/transactions"

	type Payload struct {
		Type          string   `json:"type"`
		Function      string   `json:"function"`
		TypeArguments []string `json:"type_arguments"`
		Arguments     []string `json:"arguments"`
	}

	type Signature struct {
		Type      string `json:"type"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	}

	payload := struct {
		Sender                  string `json:"sender"`
		SequenceNumber          string `json:"sequence_number"`
		MaxGasAmount            string `json:"max_gas_amount"`
		GasUnitPrice            string `json:"gas_unit_price"`
		ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`

		Payload   Payload   `json:"payload"`
		Signature Signature `json:"signature"`
	}{Sender: tx.Sender, SequenceNumber: tx.SequenceNumber, MaxGasAmount: tx.MaxGasAmount, GasUnitPrice: tx.GasUnitPrice, ExpirationTimestampSecs: tx.ExpirationTimestampSecs,
		Payload: Payload{
			Type:          "entry_function_payload",
			Function:      "0x1::aptos_coin::transfer",
			TypeArguments: []string{"string"},
			Arguments:     nil,
		},
		Signature: Signature{
			Type:      tx.Signature.Type,
			PublicKey: tx.Signature.PublicKey,
			Signature: tx.Signature.Signature,
		},
	}

	ret, err := c.doPost(route, payload)
	if err != nil {
		return "", err
	}

	res := struct {
		Hash string `json:"hash"`
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return "", err
	}

	return res.Hash, nil
}

type Tx struct {
	Type                    string `json:"type"`
	Hash                    string `json:"hash"`
	Sender                  string `json:"sender"`
	SequenceNumber          string `json:"sequence_number"`
	MaxGasAmount            string `json:"max_gas_amount"`
	GasUnitPrice            string `json:"gas_unit_price"`
	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
	Signature               struct {
		Type      string `json:"type"`
		PublicKey string `json:"public_key"`
		Signature string `json:"signature"`
	} `json:"signature"`
	Payload struct {
		Type          string        `json:"type"`
		Function      string        `json:"function"`
		TypeArguments []string      `json:"type_arguments"`
		Arguments     []interface{} `json:"arguments"`
	} `json:"payload"`
}

// GetTxByAddr get tx by address
func (c *Client) GetTxByAddr(addr string) ([]*Tx, error) {
	route := fmt.Sprintf("/v1/accounts/%s/transactions", addr)

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := make([]*Tx, 0)

	if err := json.Unmarshal(ret, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetTxByTxHash get tx by tx hash
func (c *Client) GetTxByTxHash(txHash string) (*Tx, error) {
	route := fmt.Sprintf("/v1/transactions/by_hash/%s", txHash)

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := new(Tx)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetTxByVersion get tx by version
func (c *Client) GetTxByVersion(txVersion string) (*Tx, error) {
	route := fmt.Sprintf("/v1/transactions/by_version/%s", txVersion)

	ret, err := c.doGet(route)
	if err != nil {
		return nil, err
	}

	res := new(Tx)
	if err := json.Unmarshal(ret, res); err != nil {
		return nil, err
	}

	return res, nil
}

//
// // GetBalance get balance with coin name and address
// func (c *OasisClient) GetBalance(ctx context.Context, addr string) (types.BigInt, error) {
// 	route := "/account/balance"
//
// 	type AccountIdentifier struct {
// 		Address string `json:"address"`
// 	}
//
// 	payload := struct {
// 		*baseReq          `json:"network_identifier"`
// 		AccountIdentifier `json:"account_identifier"`
// 	}{&baseReq{
// 		Blockchain: "Oasis",
// 		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
// 	}, AccountIdentifier{Address: addr}}
//
// 	ret, err := c.doPost(route, payload)
// 	if err != nil {
// 		return types.EmptyInt, err
// 	}
//
// 	fmt.Println(string(ret))
//
// 	res := struct {
// 		Balances []struct {
// 			Value string `json:"value"`
// 		}
// 	}{}
//
// 	if err := json.Unmarshal(ret, &res); err != nil {
// 		return types.EmptyInt, err
// 	}
//
// 	return types.BigFromString(res.Balances[0].Value)
// }
//
// // GetNonce from balance api
// func (c *OasisClient) GetNonce(ctx context.Context, addr string) (uint64, error) {
// 	route := "/account/balance"
//
// 	type AccountIdentifier struct {
// 		Address string `json:"address"`
// 	}
//
// 	payload := struct {
// 		*baseReq          `json:"network_identifier"`
// 		AccountIdentifier `json:"account_identifier"`
// 	}{&baseReq{
// 		Blockchain: "Oasis",
// 		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
// 	}, AccountIdentifier{Address: addr}}
//
// 	ret, err := c.doPost(route, payload)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	fmt.Println(string(ret))
//
// 	res := struct {
// 		Metadata struct {
// 			Nonce int
// 		}
// 	}{}
//
// 	if err := json.Unmarshal(ret, &res); err != nil {
// 		return 0, err
// 	}
//
// 	return uint64(res.Metadata.Nonce), nil
// }
//

//
// // GetSupportNetwork get support network
// func (c *OasisClient) GetSupportNetwork(ctx context.Context) (map[string]string, error) {
// 	route := "/network/list"
//
// 	payload := struct {
// 		*baseReq `json:"network_identifier"`
// 	}{&baseReq{
// 		Blockchain: "Oasis",
// 		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
// 	}}
//
// 	ret, err := c.doPost(route, payload)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	fmt.Println(string(ret))
//
// 	type item struct {
// 		Blockchain string `json:"blockchain"`
// 		Network    string `json:"network"`
// 	}
// 	res := struct {
// 		NetworkIdentifiers []item `json:"network_identifiers"`
// 	}{}
//
// 	if err := json.Unmarshal(ret, &res); err != nil {
// 		return nil, err
// 	}
//
// 	m := make(map[string]string)
// 	for _, v := range res.NetworkIdentifiers {
// 		m[v.Blockchain] = v.Network
// 	}
//
// 	return m, nil
// }
//
func (c *Client) doPost(route string, input interface{}) ([]byte, error) {
	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+route, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// set header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, application/x-bcs")

	// send req
	resp, err := c.httpClient.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("resp status code is not ok, code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithDetails(err, "ioutil.ReadAll fail", "resp", resp, "input", input)
	}

	return body, nil
}

func (c *Client) doGet(route string) ([]byte, error) {
	resp, err := c.httpClient.Get(c.endpoint + route)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("resp status code is not ok, code: %d", resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithDetails(err, "ioutil.ReadAll fail", "resp", resp, "route", route)
	}

	return body, nil
}
