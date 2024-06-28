package icp

import (
	"net/http"
	"net/url"
)

// Client is a client for the IC agent.
type Client struct {
	client http.Client
	config ClientConfig
}

// ClientConfig is the configuration for a client.
type ClientConfig struct {
	Host *url.URL
}

// NewClient creates a new client based on the given configuration.
func NewClient(cfg ClientConfig) Client {
	return Client{
		client: http.Client{},
		config: cfg,
	}
}
