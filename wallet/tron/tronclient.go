package tron

import (
	"errors"
	"math/big"
	"net"
	"strings"
	"sync"

	"github.com/SavourDao/savour-core/config"
	"github.com/ethereum/go-ethereum/log"
	tclient "github.com/fbsobreira/gotron-sdk/pkg/client"
)

var (
	blockNumberCacheTime int64 = 10 // seconds
)

const (
	ChainIDMain = 0x41
	//ChainIDTest = 0xa0
	ChainIDTest = 0x41
)

type tronClient struct {
	grpcClient       *tclient.GrpcClient
	chainID          byte
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

func newTronClients(conf *config.Config) ([]*tronClient, error) {
	var clients []*tronClient
	log.Info("tron client setup", "network", conf.NetWork)
	for _, rpc := range conf.Fullnode.Trx.RPCs {
		var client tronClient
		client.confirmations = conf.Fullnode.Trx.Confirmations

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
			log.Info("tronclient setup client", "ip", ipAddr)

			rpcURL = strings.Replace(rpc.RPCURL, words[0], ipAddr.String(), 1)
		}
		c := tclient.NewGrpcClient(rpcURL)
		if err := c.Start(); err != nil {
			continue
		}

		client.chainID = ChainIDTest
		if conf.NetWork == "mainnet" {
			client.chainID = ChainIDMain
		}
		client.grpcClient = c
		clients = append(clients, &client)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (t *tronClient) Close() {
	t.grpcClient.Stop()
}

func (t *tronClient) GetLatestBlockHeight() (int64, error) {
	res, err := t.grpcClient.GetNowBlock()
	if err != nil {
		return 0, err
	}
	return res.GetBlockHeader().GetRawData().GetNumber(), nil
}

func newLocalTronClient(network config.NetWorkType) *tronClient {
	var chainID byte
	switch network {
	case config.MainNet:
		chainID = ChainIDMain
	case config.TestNet:
		chainID = ChainIDTest
	case config.RegTest:
		chainID = ChainIDTest
	}
	return &tronClient{
		grpcClient:       &tclient.GrpcClient{},
		cacheBlockNumber: nil,
		chainID:          chainID,
		local:            true,
	}
}
