package config

import (
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/yaml.v2"
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
	Confirmations uint64 `yaml:"confirmations"`
}

// Fullnode define
type Fullnode struct {
	Btc Node `yaml:"btc"`
	Eth Node `yaml:"eth"`
	Trx Node `yaml:"trx"`
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
