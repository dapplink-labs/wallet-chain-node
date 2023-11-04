package egld

//import (
//	"context"
//	logger "github.com/multiversx/mx-chain-logger-go"
//	"github.com/multiversx/mx-sdk-go/blockchain"
//	"github.com/multiversx/mx-sdk-go/core"
//	"github.com/multiversx/mx-sdk-go/data"
//	"github.com/pkg/errors"
//	"github.com/savour-labs/wallet-chain-node/config"
//	"net"
//	"strings"
//	"time"
//)
//
//type EgldClient interface {
//	GetLatestBlockHeight() (int64, error)
//	GetAccountBalance(address string) *data.Account
//}
//
//type Client struct {
//	proxy     blockchain.Proxy
//	ArgsProxy blockchain.ArgsProxy
//}
//
//var log = logger.GetOrCreate("mx-client")
//
//func newElgdClient(conf *config.Config) ([]*Client, error) {
//	var clients []*Client
//	for _, rpc := range conf.Fullnode.Egld.RPCs {
//		rpcUrl := normalizeRpcUrl(rpc)
//		args := blockchain.ArgsProxy{
//			ProxyURL:            rpcUrl,
//			Client:              nil,
//			SameScState:         false,
//			ShouldBeSynced:      false,
//			FinalityCheck:       false,
//			CacheExpirationTime: time.Minute,
//			EntityType:          core.Proxy,
//		}
//		ep, err := blockchain.NewProxy(args)
//		if err != nil {
//			log.Error("error creating proxy", "error", err)
//		} else {
//			clients = append(clients, &Client{proxy: ep, ArgsProxy: args})
//		}
//	}
//	if len(clients) == 0 {
//		return nil, errors.New("No clients available")
//	}
//	return clients, nil
//}
//
//func normalizeRpcUrl(rpc *config.RPC) string {
//	rpcURL := rpc.RPCURL
//	domain := strings.TrimPrefix(rpc.RPCURL, "http://")
//	domain = strings.TrimPrefix(domain, "https://")
//	if strings.Contains(domain, ":") {
//		words := strings.Split(domain, ":")
//		ipAddr, err := net.ResolveIPAddr("ip", words[0])
//		if err != nil {
//			log.Error("resolve eth domain failed", "url", rpc.RPCURL)
//		}
//		log.Info("chain client setup client", "ip", ipAddr)
//		rpcURL = strings.Replace(rpc.RPCURL, words[0], ipAddr.String(), 1)
//	}
//	return rpcURL
//}
//
//func (c *Client) GetNetworkConfig() *data.NetworkConfig {
//	networkConfig, err := c.proxy.GetNetworkConfig(context.Background())
//	if err != nil {
//		log.Error("error getting network config", "error", err)
//		return nil
//	}
//	return networkConfig
//}
//
//func (c *Client) GetAccountBalance(address string) *data.Account {
//	addr, err := data.NewAddressFromBech32String(address)
//	if err != nil {
//		log.Error("invalid address", "error", err)
//		return nil
//	}
//	accountInfo, err := c.proxy.GetAccount(context.Background(), addr)
//	if err != nil {
//		log.Error("error retrieving account info", "error", err)
//		return nil
//	}
//	return accountInfo
//}
//
//func (c *Client) GetLatestBlockHeight() (int64, error) {
//	return 0, nil
//}
