package solana

import (
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/solscan"
)

var (
	OklinkBaseUrl = "https://www.oklink.com/"
	OklinkApiKey  = "5181d535-b68f-41cf-bbc6-25905e46b6a6"
	OkTimeout     = time.Second * 20
)

type SolScan struct {
	SolScanCli *solscan.ChainExplorerAdaptor
}

func NewSolScanClient(baseUrl, apiKey string, timeout time.Duration) (*SolScan, error) {
	solCli, err := solscan.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("New solscan client fail", "err", err)
		return nil, err
	}
	return &SolScan{SolScanCli: solCli}, err
}

func (ss *SolScan) GetTxByAddress(page, pagesize uint64, address string, action account.ActionType) (*account.TransactionResponse[account.AccountTxResponse], error) {
	request := &account.AccountTxRequest{
		PageRequest: chain.PageRequest{
			Page:  page,
			Limit: pagesize,
		},
		Action:  action,
		Address: address,
	}
	txData, err := ss.SolScanCli.GetTxByAddress(request)
	if err != nil {
		return nil, err
	}
	return txData, nil
}
