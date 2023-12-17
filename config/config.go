package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/log"
)

// Server prot
type Server struct {
	Port string `yaml:"port"`
}

// RPC connection info define
type RPC struct {
	RPCURL  string `yaml:"rpc_url"`
	RPCUser string `yaml:"rpc_user"`
	RPCPass string `yaml:"rpc_pass"`
}

type Node struct {
	RPCs          []*RPC `yaml:"rpcs"`
	TpApiUrl      string `yaml:"tp_api_url"`
	TpApiKey      string `yaml:"tp_api_key"`
	Confirmations uint64 `yaml:"confirmations"`
	ApiToken      string `yaml:"apiToken"`
}

type SolanaNode struct {
	PublicUrl        string `yaml:"public_url"`
	NetWork          string `yaml:"network"`
	NonceAccountAddr string `yaml:"NonceAccountAddr"`
	FeeAccountPriKey string `yaml:"FeeAccountPriKey"`
}

// Fullnode define
type Fullnode struct {
	Btc      Node       `yaml:"btc"`
	Eth      Node       `yaml:"eth"`
	Arbi     Node       `yaml:"arbi"`
	Op       Node       `yaml:"op"`
	Zksync   Node       `yaml:"zksync"`
	Bsc      Node       `yaml:"bsc"`
	Heco     Node       `yaml:"heco"`
	Avax     Node       `yaml:"avax"`
	Evmos    Node       `yaml:"evmos"`
	Polygon  Node       `yaml:"polygon"`
	Trx      Node       `yaml:"trx"`
	Near     Node       `yaml:"near"`
	Algo     Node       `yaml:"alog"`
	Xrp      Node       `yaml:"xrp"`
	Sol      SolanaNode `yaml:"solana"`
	Cosmos   Node       `yaml:"cosmos"`
	FileCoin Node       `yaml:"filecoin"`
	Dot      Node       `yaml:"dot"`
	Eos      Node       `yaml:"eos"`
	Oasis    Node       `yaml:"oasis"`
	Tezos    Node       `yaml:"tezos"`
	Aptos    Node       `yaml:"aptos"`
	Egld     Node       `yaml:"egld"`
	Mantle   Node       `yaml:"mantle"`
	Scroll   Node       `yaml:"scroll"`
	Base     Node       `yaml:"base"`
	Linea    Node       `yaml:"linea"`
	Ada      Node       `yaml:"ada"`
	Sui      Node       `yaml:"sui"`
	Flow     Node       `yaml:"flow"`
	Arweave  Node       `yaml:"arweave"`
}

// Config instance define
type Config struct {
	Server   Server   `yaml:"server"`
	Fullnode Fullnode `yaml:"fullnode"`
	NetWork  string   `yaml:"network"`
	Chains   []string `yaml:"chains"`
}

type NetWorkType int

const (
	MainNet NetWorkType = iota
	TestNet
	RegTest
)

// Setup init config
func New(path string) (*Config, error) {
	// config global config instance
	var config = new(Config)
	h := log.StreamHandler(os.Stdout, log.TerminalFormat(true))
	log.Root().SetHandler(h)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

const UnsupportedChain = "Unsupport chain"
const UnsupportedOperation = UnsupportedChain
