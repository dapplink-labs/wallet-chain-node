package tron

import (
	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"math/big"
	"net"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/log"
	tclient "github.com/fbsobreira/gotron-sdk/pkg/client"

	"github.com/savour-labs/wallet-chain-node/config"
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

// newTronClient init the tron client
func newTronClients(conf *config.Config) ([]*tronClient, error) {
	var clients []*tronClient
	chainConfig := params.SepoliaChainConfig
	if conf.NetWork == "mainnet" {
		chainConfig = params.MainnetChainConfig
	} else if conf.NetWork == "regtest" {
		chainConfig = params.AllCliqueProtocolChanges
	}
	log.Info("tron client setup", "chain_id", chainConfig.ChainID.Int64(), "network", conf.NetWork)

	for _, rpc := range conf.Fullnode.Tron.RPCs {

		var client tronClient
		client.confirmations = conf.Fullnode.Tron.Confirmations

		rpcURL := strings.Split(rpc.RPCURL, "=")[0]

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
		if err := c.Start(grpc.WithTransportCredentials(insecure.NewCredentials())); err != nil {
			log.Error("tron client start failed", "url", rpc.RPCURL, "err", err)
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
