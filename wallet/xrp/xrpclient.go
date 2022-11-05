package xrp

import (
	"encoding/json"
	"github.com/SavourDao/savour-hd/config"
	xrprpc "github.com/SavourDao/savour-hd/wallet/xrp/rpc"
	"github.com/SavourDao/savour-hd/wallet/xrp/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"sync"
)

type Client struct {
	RpcClient        xrprpc.RpcClient
	nodeConfig       config.Node
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

func (c *Client) GetLatestBlockHeight() (int64, error) {
	var res types.GetBlockResult
	e := c.RpcClient.Request(`{"method": "ledger_current","params": []}`, &res)
	if e != nil {
		log.Error("get transaction receipt error", "err", e)
		return 0, e
	}
	return int64(res.Result.LedgerCurrentIndex), nil
}

func (c *Client) GetBalance(address string) (string, error) {
	var res types.GetBalanceRes
	e := c.RpcClient.RpcRequest("account_info", []any{
		types.GetBalance{
			Strict:      true,
			Account:     address,
			LedgerIndex: "current",
			Queue:       true,
		},
	}, &res)
	if e != nil {
		log.Error("get transaction receipt error", "err", e)
		return "0", e
	}
	return res.Result.AccountData.Balance, nil
}

func (c *Client) GetTxsByAddress(address string) ([]types.Transaction, error) {
	var res types.GetAccountResult
	var param = types.RequestBody{
		Params: []any{
			types.GetAccountTx{
				Account: address,
				Binary:  false,
				Forward: true,
			},
		},
		Method: "account_tx",
	}
	jsonBody, _ := json.Marshal(param)
	e := c.RpcClient.Request(string(jsonBody), &res)
	if e != nil {
		log.Error("get GetTxsByAddress error", "err", e)
		return nil, e
	}
	var txs = make([]types.Transaction, 0)
	for _, item := range res.Result.Transactions {
		tx := types.Transaction{
			Amount:      item.Tx.Amount,
			Fee:         item.Tx.Fee,
			To:          item.Tx.Destination,
			From:        item.Tx.Account,
			Hash:        item.Tx.Hash,
			BlockHeight: string(rune(item.Tx.InLedger)),
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

func (c *Client) GetTxByHash(hash string) (*types.Transaction, error) {
	var res types.GetTxByHashResult
	e := c.RpcClient.RpcRequest("tx", []any{
		types.GetTx{
			Transaction: hash,
		},
	}, &res)
	if e != nil {
		log.Error("get GetTxByHash receipt error", "err", e)
		return nil, e
	}
	tx := types.Transaction{
		Amount:      res.Result.Amount.Value,
		Fee:         res.Result.Fee,
		To:          res.Result.Destination,
		From:        res.Result.Account,
		Hash:        res.Result.Hash,
		BlockHeight: string(rune(res.Result.InLedger)),
	}
	return &tx, nil
}

func (c *Client) SendTx(signedTx string) (string, error) {
	var res types.SendTxResult
	e := c.RpcClient.RpcRequest("submit", []any{
		types.SendTx{
			TxBlob: signedTx,
		},
	}, &res)
	if e != nil {
		log.Error("get SendTx error", "err", e)
		return "", e
	}
	return res.Result.TxJson.Hash, nil
}

func newClients(conf *config.Config) ([]*Client, error) {
	var clients []*Client
	for _, rpc := range conf.Fullnode.Xrp.RPCs {
		clients = append(clients, &Client{
			RpcClient: xrprpc.RpcClient{
				URL: rpc.RPCURL,
			},
			nodeConfig: conf.Fullnode.Xrp,
		})
	}
	return clients, nil
}
