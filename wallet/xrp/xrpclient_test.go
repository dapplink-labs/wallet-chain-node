package xrp

import (
	"fmt"
	"testing"

	"github.com/savour-labs/wallet-hd-chain/config"
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

func TestGetTxsByAddress(t *testing.T) {
	client := newTestClient()
	height, _ := client.GetTxsByAddress("rHZQVErv6UZTtkM4NVexS1Z5dTvuek2RbX")
	fmt.Println(height)
}

func TestGetTxByHash(t *testing.T) {
	client := newTestClient()
	height, _ := client.GetTxByHash("rHZQVErv6UZTtkM4NVexS1Z5dTvuek2RbX")
	fmt.Println(height)
}

func TestSendTx(t *testing.T) {
	client := newTestClient()
	height, _ := client.SendTx("xxx")
	fmt.Println(height)
}
