package near

import (
	"fmt"
	"github.com/SavourDao/savour-hd/config"
	"testing"
)

func newTestClient() *nearClient {

	client, _ := newNearClients(&config.Config{Fullnode: config.Fullnode{
		Near: config.Node{
			RPCs: []*config.RPC{
				{
					RPCURL: "https://rpc.mainnet.near.org",
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

func TestRequest(t *testing.T) {
	client := newTestClient()
	client.Request()
}
