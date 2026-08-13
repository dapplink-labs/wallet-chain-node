package ethereum

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/rpc"
	"math/big"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"

	"github.com/savour-labs/wallet-chain-node/config"
)

var (
	blockNumberCacheTime int64 = 10 // seconds
)

type ethClient struct {
	Client
	chainConfig      *params.ChainConfig
	cacheBlockNumber *big.Int
	cacheTime        int64
	rw               sync.RWMutex
	confirmations    uint64
	local            bool
}

type Client interface {
	bind.ContractBackend
	BalanceAt(context.Context, common.Address, *big.Int) (*big.Int, error)
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
	NonceAt(context.Context, common.Address, *big.Int) (uint64, error)
	Client() *rpc.Client
}

// newEthClient init the eth client
func newEthClients(conf *config.Config) ([]*ethClient, error) {
	chainConfig := params.SepoliaChainConfig
	if conf.NetWork == "mainnet" {
		chainConfig = params.MainnetChainConfig
	} else if conf.NetWork == "regtest" {
		chainConfig = params.AllCliqueProtocolChanges
	}
	log.Info("eth client setup", "chain_id", chainConfig.ChainID.Int64(), "network", conf.NetWork)

	var clients []*ethClient
	for _, rpc := range conf.Fullnode.Eth.RPCs {
		client := &ethClient{
			chainConfig:   chainConfig,
			confirmations: conf.Fullnode.Eth.Confirmations,
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
		client.Client, err = ethclient.Dial(rpcURL)
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

func newLocalEthClient(network config.NetWorkType) *ethClient {
	var para *params.ChainConfig
	switch network {
	case config.MainNet:
		para = params.MainnetChainConfig
	case config.TestNet:
		para = params.SepoliaChainConfig
	case config.RegTest:
		para = params.AllCliqueProtocolChanges
	default:
		panic("unsupported network type")
	}
	return &ethClient{
		Client:           &ethclient.Client{},
		chainConfig:      para,
		cacheBlockNumber: nil,
		local:            true,
	}
}

func (client *ethClient) blockNumber() *big.Int {
	now := time.Now().Unix()
	client.rw.RLock()
	if now-client.cacheTime < blockNumberCacheTime {
		number := client.cacheBlockNumber
		client.rw.RUnlock()
		return number
	}
	client.rw.RUnlock()

	client.rw.Lock()
	defer client.rw.Unlock()
	if now-client.cacheTime < blockNumberCacheTime {
		return client.cacheBlockNumber
	}
	latestBlock, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Error("get BlockByNumber failed", "error", err)
		return nil
	}
	client.cacheBlockNumber = latestBlock.Number()
	client.cacheTime = now
	return client.cacheBlockNumber
}

func (client *ethClient) isContractAddress(address common.Address) bool {
	code, err := client.CodeAt(context.Background(), address, nil)
	return err == nil && len(code) > 0
}

func (client *ethClient) GetLatestBlockHeight() (int64, error) {
	number, err := client.BlockByNumber(context.TODO(), nil)
	if err != nil {
		return 0, err
	}
	return number.Number().Int64(), err
}
