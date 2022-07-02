package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database  *Database  `yaml:"database"`
	RpcServer *RpcServer `yaml:"rpc_server"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type RpcServer struct {
	Ip   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func LoadConfigFile(filePath string, cfg *Config) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return yaml.NewDecoder(file).Decode(cfg)
}
