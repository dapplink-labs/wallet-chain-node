package xrp

import (
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

func newClients(conf *config.Config) ([]*Client, error) {
	var clients []*Client
	for _, rpc := range conf.Fullnode.Xrp.RPCs {
		clients = append(clients, &Client{
			RpcClient: xrprpc.RpcClient{
				URL: rpc.RPCURL,
			},
			nodeConfig: conf.Fullnode.Near,
		})
	}
	return clients, nil
}
