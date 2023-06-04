package eosio

import (
	eos "github.com/eoscanada/eos-go"
)

type Client struct {
	client *eos.API
}

func newClient(url string) (*Client, error) {

	return &Client{
		client: eos.New(url),
	}, nil
}
