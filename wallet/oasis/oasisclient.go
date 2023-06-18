package oasis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"emperror.dev/errors"
	"github.com/filecoin-project/lotus/chain/types"
)

type TransactionIdentifier struct {
	Hash string `json:"hash"`
}

type Transaction struct {
	TransactionIdentifier *TransactionIdentifier `json:"transaction_identifier"`
	Metadata              interface{}            `json:"metadata"`
}

type OasisClient struct {
	cflb        string
	endpoint    string
	oasisClient *http.Client
}

func NewOasisClient(cflb, endpoint string) (*OasisClient, error) {
	return &OasisClient{
		cflb:        cflb,
		endpoint:    endpoint,
		oasisClient: new(http.Client),
	}, nil
}

type baseReq struct {
	Blockchain string `json:"blockchain"`
	Network    string `json:"network"`
}

// GetBalance get balance with coin name and address
func (c *OasisClient) GetBalance(ctx context.Context, addr string) (types.BigInt, error) {
	route := "/account/balance"

	type AccountIdentifier struct {
		Address string `json:"address"`
	}

	payload := struct {
		*baseReq          `json:"network_identifier"`
		AccountIdentifier `json:"account_identifier"`
	}{&baseReq{
		Blockchain: "Oasis",
		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
	}, AccountIdentifier{Address: addr}}

	ret, err := c.doPost(route, payload)
	if err != nil {
		return types.EmptyInt, err
	}

	fmt.Println(string(ret))

	res := struct {
		Balances []struct {
			Value string `json:"value"`
		}
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return types.EmptyInt, err
	}

	return types.BigFromString(res.Balances[0].Value)
}

// GetNonce from balance api
func (c *OasisClient) GetNonce(ctx context.Context, addr string) (uint64, error) {
	route := "/account/balance"

	type AccountIdentifier struct {
		Address string `json:"address"`
	}

	payload := struct {
		*baseReq          `json:"network_identifier"`
		AccountIdentifier `json:"account_identifier"`
	}{&baseReq{
		Blockchain: "Oasis",
		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
	}, AccountIdentifier{Address: addr}}

	ret, err := c.doPost(route, payload)
	if err != nil {
		return 0, err
	}

	fmt.Println(string(ret))

	res := struct {
		Metadata struct {
			Nonce int
		}
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return 0, err
	}

	return uint64(res.Metadata.Nonce), nil
}

// BroadcastTx send tx to the network
func (c *OasisClient) BroadcastTx(ctx context.Context, txByte []byte) (string, error) {
	route := "/construction/submit"

	payload := struct {
		*baseReq          `json:"network_identifier"`
		SignedTransaction string `json:"signed_transaction"`
	}{&baseReq{
		Blockchain: "Oasis",
		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
	}, string(txByte)}

	ret, err := c.doPost(route, payload)
	if err != nil {
		return "", err
	}

	res := struct {
		TransactionIdentifier struct {
			Hash string `json:"hash"`
		}
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return "", err
	}

	return res.TransactionIdentifier.Hash, nil
}

// GetSupportNetwork get support network
func (c *OasisClient) GetSupportNetwork(ctx context.Context) (map[string]string, error) {
	route := "/network/list"

	payload := struct {
		*baseReq `json:"network_identifier"`
	}{&baseReq{
		Blockchain: "Oasis",
		Network:    "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535",
	}}

	ret, err := c.doPost(route, payload)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(ret))

	type item struct {
		Blockchain string `json:"blockchain"`
		Network    string `json:"network"`
	}
	res := struct {
		NetworkIdentifiers []item `json:"network_identifiers"`
	}{}

	if err := json.Unmarshal(ret, &res); err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, v := range res.NetworkIdentifiers {
		m[v.Blockchain] = v.Network
	}

	return m, nil
}

func (c *OasisClient) doPost(route string, input interface{}) ([]byte, error) {
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
	req.Header.Set("Cookie", fmt.Sprintf("__cflb=%s", c.cflb))

	// send req
	resp, err := c.oasisClient.Do(req)
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
