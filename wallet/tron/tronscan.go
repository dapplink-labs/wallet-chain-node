package tron

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
)

type TronScan struct {
	TronScanCli *oklink.ChainExplorerAdaptor
}

func NewTronScanClient(baseUrl, apiKey string, timeout time.Duration) (*TronScan, error) {
	tronCli, err := oklink.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("new tron scan client fail", "err", err)
		return nil, err
	}
	return &TronScan{TronScanCli: tronCli}, err
}

func (ts *TronScan) GetTxByAddress(page, pagesize uint64, address string) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		ChainShortName: ChainName,
		ExplorerName:   oklink.ChainExplorerName,
		Action:         account.OkLinkActionNormal,
		Address:        address,
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
	}
	txData, err := ts.TronScanCli.GetTxByAddress(request)
	if err != nil {
		return nil, err
	}
	return txData, nil
}
