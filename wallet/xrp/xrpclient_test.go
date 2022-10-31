package xrp

import (
	"fmt"
	"github.com/SavourDao/savour-hd/config"
	"testing"
)

func newTestClient() *Client {

	client, _ := newClients(&config.Config{Fullnode: config.Fullnode{
		Xrp: config.Node{
			RPCs: []*config.RPC{
				{
					RPCURL: "https://xrplcluster.com",
				},
			},
		},
	},
	})
	return client[0]
}

func TestGetLatestBlockHeight(t *testing.T) {
	client := newTestClient()
	height, _ := client.GetLatestBlockHeight()
	fmt.Println(height)
}

func TestGetBalance(t *testing.T) {
	client := newTestClient()
	height, _ := client.GetBalance("rEtuF2YJzoRQCAaYAtBumpm4UWQMiV4MXJ")
	fmt.Println(height)
}
