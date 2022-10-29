package near

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/SavourDao/savour-hd/config"
	"github.com/SavourDao/savour-hd/wallet/near/account"
	"github.com/SavourDao/savour-hd/wallet/near/keys"
	"github.com/SavourDao/savour-hd/wallet/near/transaction"
	"github.com/SavourDao/savour-hd/wallet/near/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/golang-module/dongle"
	"github.com/portto/solana-go-sdk/client"

	"github.com/shopspring/decimal"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"sync"
)

type NearClient struct {
	RpcClient        *rpc.Client
	Client           client.Client
	nodeConfig       config.Node
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

func (c *NearClient) GetBlock() (types.GetBlockResult, error) {
	jsonStr := []byte(`{"jsonrpc": "2.0", "method": "block", "params": {"finality": "final"},"id": 1}`)
	url := c.nodeConfig.RPCs[0].RPCURL
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.GetBlockResult{}
	_ = json.Unmarshal(body, &res)
	return res, nil
}

func (c *NearClient) GetLatestBlockHeight() (int64, error) {
	res, e := c.GetBlock()
	if e != nil {
		return 0, nil
	}
	return int64(res.Result.Header.Height), nil
}

func (c *NearClient) GetBalance(address string) (string, error) {
	reqParams := types.RPCRequest{
		Jsonrpc: "2.0",
		Method:  "query",
		Params: types.GetBalanceParam{
			Finality:    "final",
			RequestType: "view_account",
			AccountID:   address,
		},
		ID: 1,
	}
	jsonStr, _ := json.Marshal(reqParams)
	url := c.nodeConfig.RPCs[0].RPCURL
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	res := types.GetBalanceRes{}
	_ = json.Unmarshal(body, &res)
	var amount, _ = strconv.ParseFloat(res.Result.Amount, 64)
	decima := math.Pow(10, 24)
	d2 := decimal.NewFromFloat(amount).Div(decimal.NewFromFloat(decima))
	return d2.String(), nil
}

func (c *NearClient) GetTx(address string, page int, size int) ([]Transaction, error) {

	var txs = make([]Transaction, 5)
	sqlStr := "select * from transactions where receiver_account_id = $1 order by block_timestamp desc limit $2 offset $3"
	rows, err := db.Query(sqlStr, address, page, size)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tx Transaction
		err := rows.Scan(
			&tx.TransactionHash, &tx.IncludedInBlockHash, &tx.IncludedInChunkHash, &tx.IndexInChunk, &tx.BlockTimestamp, &tx.SignerAccountId, &tx.SignerPublicKey,
			&tx.Nonce, &tx.ReceiverAccountId, &tx.Signature, &tx.Status, &tx.ConvertedIntoReceiptId, &tx.ReceiptConversionGasBurnt, &tx.ReceiptConversionTokensBurnt,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		txs = append(txs, tx)
		fmt.Printf("student=%+v\n", tx)
	}
	return nil, nil
}

func (c *NearClient) GetAccount() (string, string, error) {
	publicKey, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}
	pri := dongle.Encode.FromBytes(priv).ByBase58().ToString()
	address := hex.EncodeToString(publicKey)
	return pri, address, nil
}

func (c *NearClient) GetTxByHash(hash string) (*Transaction, error) {
	sqlStr := "select * from transactions where transaction_hash = $1 limit 1"
	rows, err := db.Query(sqlStr, hash)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	var tx Transaction
	if rows.Next() {
		err := rows.Scan(
			&tx.TransactionHash, &tx.IncludedInBlockHash, &tx.IncludedInChunkHash, &tx.IndexInChunk, &tx.BlockTimestamp, &tx.SignerAccountId, &tx.SignerPublicKey,
			&tx.Nonce, &tx.ReceiverAccountId, &tx.Signature, &tx.Status, &tx.ConvertedIntoReceiptId, &tx.ReceiptConversionGasBurnt, &tx.ReceiptConversionTokensBurnt,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
		}
		fmt.Printf("student=%+v\n", tx)
	}
	fmt.Println(tx)
	return &tx, nil

}

func (c *NearClient) RpcRequest(method string, params any, result interface{}) error {
	reqParams := types.RPCRequest{
		Jsonrpc: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	jsonStr, _ := json.Marshal(reqParams)
	url := c.nodeConfig.RPCs[0].RPCURL
	req, e := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, e := httpClient.Do(req)
	if e != nil {
		return e
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &result)
	fmt.Println(string(body))
	return nil
}

func (c *NearClient) getAccessKey(pub string, from string, res *types.GetAccessKeyResponse) error {
	param := types.GetAccessKeyRequest{
		Finality:    "final",
		RequestType: "view_access_key",
		AccountID:   from,
		PublicKey:   pub,
	}
	c.RpcRequest("query", param, &res)
	return nil
}

func (c *NearClient) SendTx(pri string, from string, to string, amount string) (string, error) {

	signer, err := keys.NewKeyPairFromString("ed25519:" + pri)
	if err != nil {
		fmt.Println("NewKeyPairFromString error")
		return "", err
	}
	key := signer.GetPublicKey()
	pub, _ := key.ToString()

	accessKeyResponse := types.GetAccessKeyResponse{}
	c.getAccessKey(pub, from, &accessKeyResponse)
	block, e := c.GetBlock()
	if e != nil {
		fmt.Println("GetBlock error")
		return "", e
	}
	acc := account.NewAccount(&types.Config{
		Signer:    signer,
		NetworkID: "mainnet",
	}, from)
	amt, _ := (&big.Int{}).SetString(amount, 10)
	sendAction := transaction.TransferAction(*amt)
	sendTransactionResult, err := acc.SignAndSendTransaction(context.Background(), to, accessKeyResponse, block.Result.Header.Hash, sendAction)
	fmt.Println(sendTransactionResult)
	if err != nil {
		return "", err
	}
	return sendTransactionResult.Result.Transaction.Hash, nil
}

func newNearClients(conf *config.Config) ([]*NearClient, error) {
	var clients []*NearClient
	rpcClient, e := rpc.DialContext(context.Background(), "https://rpc.testnet.near.org")
	if e != nil {
		return nil, e
	}
	clients = append(clients, &NearClient{
		RpcClient:  rpcClient,
		nodeConfig: conf.Fullnode.Near,
	})
	return clients, nil
}
