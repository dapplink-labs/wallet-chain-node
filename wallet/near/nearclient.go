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
	"reflect"
	"strings"

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
	return c.ToNearUnit(res.Result.Amount), nil
}

func (c *NearClient) ToNearUnit(amountStr string) string {
	var amount, _ = strconv.ParseFloat(amountStr, 64)
	decima := math.Pow(10, 24)
	d2 := decimal.NewFromFloat(amount).Div(decimal.NewFromFloat(decima))
	return d2.String()
}

func (c *NearClient) GetTx(address string, page int, size int) ([]BlockTransaction, error) {

	var txs = make([]BlockTransaction, 0)
	sqlStr := "select * from transactions where receiver_account_id = $1 order by block_timestamp desc limit $2 offset $3"
	rows, err := db.Query(sqlStr, address, page, size)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	var hashList []string
	var blockHashList []string
	for rows.Next() {
		var tx BlockTransaction
		err := rows.Scan(
			&tx.TransactionHash, &tx.IncludedInBlockHash, &tx.IncludedInChunkHash, &tx.IndexInChunk, &tx.BlockTimestamp, &tx.SignerAccountId, &tx.SignerPublicKey,
			&tx.Nonce, &tx.ReceiverAccountId, &tx.Signature, &tx.Status, &tx.ConvertedIntoReceiptId, &tx.ReceiptConversionGasBurnt, &tx.ReceiptConversionTokensBurnt,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
		hashList = append(hashList, tx.TransactionHash)
		blockHashList = append(blockHashList, tx.IncludedInBlockHash)
		txs = append(txs, tx)
		fmt.Printf("student=%+v\n", tx)
	}
	actions, e := c.GetActionByHash(hashList)
	if e != nil {
		fmt.Printf("GetActionByHash failed, err:%v\n", e)
		return nil, e
	}
	blocks, e := c.GetBlockListByHash(blockHashList)
	if e != nil {
		fmt.Printf("GetBlockListByHash failed, err:%v\n", e)
		return nil, e
	}
	actionMap := ListToMap(actions, "TransactionHash")
	blockMap := ListToMap(blocks, "BlockHash")

	for i := 0; i < len(txs); i++ {
		action := actionMap[txs[i].TransactionHash]
		if action != nil {
			a := action.(TransactionAction)
			var args = &TransactionActionArgs{}
			json.Unmarshal([]byte(a.Args), args)
			txs[i].Amount = args.Deposit
		}
		block := blockMap[txs[i].IncludedInBlockHash]
		if block != nil {
			a := block.(Block)
			txs[i].BlockHeight = a.BlockHeight
		}
	}
	return txs, nil
}

func ListToMap(list interface{}, key string) map[string]interface{} {
	res := make(map[string]interface{})
	arr := ToSlice(list)
	for _, row := range arr {
		immutable := reflect.ValueOf(row)
		val := immutable.FieldByName(key).String()
		res[val] = row
	}
	return res
}

func ToSlice(arr interface{}) []interface{} {
	ret := make([]interface{}, 0)
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		ret = append(ret, arr)
		return ret
	}
	l := v.Len()
	for i := 0; i < l; i++ {
		ret = append(ret, v.Index(i).Interface())
	}
	return ret
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

func (c *NearClient) GetTxByHash(hash string) (*BlockTransaction, error) {
	sqlStr := "select * from transactions where transaction_hash = $1 limit 1"
	rows, err := db.Query(sqlStr, hash)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	var tx BlockTransaction
	if rows.Next() {
		err := rows.Scan(
			&tx.TransactionHash, &tx.IncludedInBlockHash, &tx.IncludedInChunkHash, &tx.IndexInChunk, &tx.BlockTimestamp, &tx.SignerAccountId, &tx.SignerPublicKey,
			&tx.Nonce, &tx.ReceiverAccountId, &tx.Signature, &tx.Status, &tx.ConvertedIntoReceiptId, &tx.ReceiptConversionGasBurnt, &tx.ReceiptConversionTokensBurnt,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			return nil, err
		}
	}
	blocks, _ := c.GetBlockListByHash([]string{tx.IncludedInBlockHash})

	if len(blocks) > 0 {
		tx.BlockHeight = blocks[0].BlockHeight
	}
	actions, _ := c.GetActionByHash([]string{tx.TransactionHash})
	if len(actions) > 0 {
		var args = &TransactionActionArgs{}
		json.Unmarshal([]byte(actions[0].Args), args)
		tx.Amount = args.Deposit
	}
	return &tx, nil

}

func (c *NearClient) GetBlockListByHash(hashList []string) ([]Block, error) {
	var blocks = make([]Block, 0)

	sqlStr := `select * from blocks where block_hash in ('` + strings.Join(hashList, `","`) + `')`
	fmt.Printf(sqlStr)
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var block Block

		err := rows.Scan(
			&block.BlockHeight, &block.BlockHash, &block.PrevBlockHash, &block.BlockTimestamp, &block.TotalSupply, &block.GasPrice, &block.AuthorAccountId,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		blocks = append(blocks, block)
		fmt.Printf("block=%+v\n", block)
	}
	return blocks, nil
}

func (c *NearClient) GetActionByHash(hashList []string) ([]TransactionAction, error) {
	var actions = make([]TransactionAction, 0)

	sqlStr := `select * from transaction_actions where action_kind = 'TRANSFER' AND transaction_hash in ('` + strings.Join(hashList, `","`) + `')`
	fmt.Printf(sqlStr)
	rows, err := db.Query(sqlStr)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var action TransactionAction

		err := rows.Scan(
			&action.TransactionHash, &action.IndexInTransaction, &action.ActionKind, &action.Args,
		)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		actions = append(actions, action)
		fmt.Printf("action=%+v\n", action)
	}
	return actions, nil
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

func (c *NearClient) SignAndSendTx(pri string, from string, to string, amount string) (string, error) {

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

func (c *NearClient) GetNonce(pub string, from string) (string, error) {
	accessKeyResponse := types.GetAccessKeyResponse{}
	e := c.getAccessKey(pub, from, &accessKeyResponse)
	if e != nil {
		return "", e
	}
	return strconv.FormatInt(accessKeyResponse.Result.Nonce, 10), nil
}

func (c *NearClient) SendSignedTx(signedTx string) (string, error) {
	var res types.SendTxResult
	e := types.DoRpcRequest("broadcast_tx_commit", [1]string{signedTx}, &res)
	if e != nil {
		return "", e
	}
	if res.Error != nil {
		return "", fmt.Errorf("SendSignedTx error")
	}
	return res.Result.Transaction.Hash, nil
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
