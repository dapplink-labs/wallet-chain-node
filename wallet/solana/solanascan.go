package solana

import (
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/explorer/oklink"
)

var (
	OklinkBaseUrl = "https://www.oklink.com/"
	OklinkApiKey  = "5181d535-b68f-41cf-bbc6-25905e46b6a6"
	OkTimeout     = time.Second * 20
)

type SolanaExplorer struct {
	OkLinkCli *oklink.ChainExplorerAdaptor
}

func NewScanClient(baseUrl, apiKey string, timeout time.Duration) (*SolanaExplorer, error) {
	oklinkCli, err := oklink.NewChainExplorerAdaptor(apiKey, baseUrl, false, time.Duration(timeout))
	if err != nil {
		log.Error("Mock oklink client fail", "err", err)
		return nil, err
	}
	return &SolanaExplorer{OkLinkCli: oklinkCli}, err
}

func (se *SolanaExplorer) GetTxByAddress() {
	request := &account.AccountTxRequest{
		ChainShortName: "ETH",
		ExplorerName:   oklink.ChainExplorerName,
		Action:         account.OkLinkActionNormal,
		Address:        "0x4838B106FCe9647Bdf1E7877BF73cE8B0BAD5f97",
		IsFromOrTo:     string(account.From),
		PageRequest: chain.PageRequest{
			Page:  1,
			Limit: 10,
		},
	}
	se.OkLinkCli.GetTxByAddress(request)
}
