package algo

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/ethereum/go-ethereum/log"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/common"
	altype "github.com/algorand/go-algorand-sdk/types"
	"github.com/savour-labs/wallet-chain-node/config"
)

type AlgoClient interface {
	GetLatestBlockHeight() (int64, error)
	GetAccountBalance(address string) *algod.AccountInformation
	GetTransactionParams(ctx context.Context) (altype.SuggestedParams, error)
	getTransactionByHash(round uint64, txid string) *algod.GetTransactionProof
}

type Client struct {
	client        *algod.Client
	confirmations uint64
}

func newAlgoClients(conf *config.Config) ([]*Client, error) {
	var clients []*Client
	for _, rpc := range conf.Fullnode.Algo.RPCs {
		client := &Client{
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
		client.client, err = algod.MakeClient(rpcURL, conf.Fullnode.Algo.ApiToken)
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

func newLocalAlgoClient(network config.NetWorkType) *Client {
	return &Client{
		client: &algod.Client{},
	}
}

func GetLatestBlockHeight() (int64, error) {
	return 0, nil
}

func (c Client) GetAccountBalance(address string) *algod.AccountInformation {
	return c.client.AccountInformation(address)
}

func (c Client) GetTransactionParams(ctx context.Context) (altype.SuggestedParams, error) {
	h := common.Header{}
	return c.client.SuggestedParams().Do(ctx, &h)
}

func (c Client) SendRawtransaction(rawtxn []byte) *algod.SendRawTransaction {
	return c.client.SendRawTransaction(rawtxn)
}

func (c Client) getTransactionByHash(round uint64, txid string) *algod.GetTransactionProof {
	return c.client.GetTransactionProof(round, txid)
}
