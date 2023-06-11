package algo

import (
	"context"
	"errors"
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	altype "github.com/algorand/go-algorand-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/savour-labs/wallet-hd-chain/config"
	"math/big"
	"net"
	"strings"
)

type algoClient struct {
	Client        *algod.Client
	confirmations uint64
}

type Client interface {
	bind.ContractBackend

	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
}

func newAlgoClients(conf *config.Config) ([]*algoClient, error) {
	var clients []*algoClient
	for _, rpc := range conf.Fullnode.Algo.RPCs {
		client := &algoClient{
			confirmations: conf.Fullnode.Algo.Confirmations,
		}
		rpcURL := rpc.RPCURL
		domain := strings.TrimPrefix(rpc.RPCURL, "http://")
		domain = strings.TrimPrefix(domain, "https://")
		if strings.Contains(domain, ":") {
			words := strings.Split(domain, ":")
			ipAddr, err := net.ResolveIPAddr("ip", words[0])
			if err != nil {
				log.Error("resolve eth domain failed", "url", rpc.RPCURL)
				continue
			}
			log.Info("ethclient setup client", "ip", ipAddr)
			rpcURL = strings.Replace(rpc.RPCURL, words[0], ipAddr.String(), 1)
		}
		var err error
		client.Client, err = algod.MakeClient(rpcURL, conf.Fullnode.Algo.ApiToken)
		if err != nil {
			log.Error("ethclient dial failed", "err", err)
			continue
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func newLocalAlgoClient(network config.NetWorkType) *algoClient {
	return &algoClient{
		Client: &algod.Client{},
	}
}

func (a algoClient) GetLatestBlockHeight() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (a algoClient) GetAccountBalance(address string) *algod.AccountInformation {
	return a.Client.AccountInformation(address)
}

func (a algoClient) GetTransactionParams(ctx context.Context) (altype.SuggestedParams, error) {
	h := common.Header{}
	return a.Client.SuggestedParams().Do(ctx, &h)
}

func (a algoClient) Send_transaction(rawtxn []byte) *algod.SendRawTransaction {
	return a.Client.SendRawTransaction(rawtxn)

}
