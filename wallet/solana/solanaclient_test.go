package solana

import (
	"fmt"
	"github.com/dapplink-labs/chain-explorer-api/common/account"
	"github.com/dapplink-labs/chain-explorer-api/common/chain"
	"github.com/dapplink-labs/chain-explorer-api/explorer/solscan"
	"testing"
	"time"

	"github.com/savour-labs/wallet-chain-node/config"
)

func newTestClient() *SolanaClient {
	var rpcList []*config.RPC
	rpcc := &config.RPC{
		RPCURL: "https://docs-demo.solana-mainnet.quiknode.pro",
	}
	rpcList = append(rpcList, rpcc)
	client, _ := NewSolanaClients(&config.Config{
		Fullnode: config.Fullnode{
			Sol: config.SolanaNode{
				RPCs:               rpcList,
				NetWork:            "mainnet",
				SolScanBaseTimeout: 20 * time.Second,
				SolScanBaseUrl:     "https://pro-api.solscan.io",
				SolScanApiKey:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MjQwNjIxMzk5MDYsImVtYWlsIjoiemFja2d1by5ndW9AZ21haWwuY29tIiwiYWN0aW9uIjoidG9rZW4tYXBpIiwiYXBpVmVyc2lvbiI6InYxIiwiaWF0IjoxNzI0MDYyMTM5fQ.EaWDC25lyGNx_LqRL5sAYYKLMbq10brnexKnAz9C3UY",
			},
		},
	})
	return client
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("4Y19AR3cQ76UmLPeEYsvwkUXaiS8GEfivyggcYSmL4M6")
	fmt.Printf("%s", balance)
}

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := newTestClient()
	hash, err := client.GetTxByHash("2VPumjVuk6wRymhbkx2na4EDQNgfaAz63HrKye5suHEZf6y2R8grv4QpA5VCuQw5CtcPWMwfau2JEaXNbtMXesA8")
	if err != nil {
		return
	}
	fmt.Println("Hash===", hash.Hash)
	fmt.Println("Type===", hash.From)
	fmt.Println("Type===", hash.To)
	fmt.Println("Type===", hash.Type)
	fmt.Println("Value===", hash.Value)
	fmt.Println("Value===", hash.Fee)
	fmt.Println("Status===", hash.Status)
	fmt.Println("Height===", hash.Height)
}

func TestSolanaClient_GetTransferHistory(t *testing.T) {
	sol, err := solscan.NewChainExplorerAdaptor("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOjE3MjQwNjIxMzk5MDYsImVtYWlsIjoiemFja2d1by5ndW9AZ21haWwuY29tIiwiYWN0aW9uIjoidG9rZW4tYXBpIiwiYXBpVmVyc2lvbiI6InYxIiwiaWF0IjoxNzI0MDYyMTM5fQ.EaWDC25lyGNx_LqRL5sAYYKLMbq10brnexKnAz9C3UY", "https://pro-api.solscan.io", true, time.Second*20)
	if err == nil {
		request := account.AccountTxRequest{
			PageRequest: chain.PageRequest{
				Page:  1,
				Limit: 50,
			},
			Address: "4Y19AR3cQ76UmLPeEYsvwkUXaiS8GEfivyggcYSmL4M6",
		}
		resp, _ := sol.GetTxByAddress(&request)
		fmt.Println(resp)
	} else {
		fmt.Println(err)
	}
}

func TestSolanaClient_GetAccount(t *testing.T) {
	client := newTestClient()
	address, pri, _ := client.GetAccount()
	fmt.Println(address)
	fmt.Println(pri)
}

func TestSolanaClient_RequestAirdrop(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("4Y19AR3cQ76UmLPeEYsvwkUXaiS8GEfivyggcYSmL4M6")
	fmt.Println(balance)
}

func TestSolanaClient_SendTx(t *testing.T) {
	client := newTestClient()
	tx, err := client.SendTx("")
	if err != nil {
		return
	}
	fmt.Println(tx)
}

func TestSolanaClient_GetNonce(t *testing.T) {
	client := newTestClient()
	nonce, err := client.GetNonce("")
	if err != nil {
		return
	}
	fmt.Println("nonce===", nonce)
}

func TestSolanaClient_GetMinRent(t *testing.T) {
	client := newTestClient()
	rent, err := client.GetMinRent()
	if err != nil {
		return
	}
	fmt.Println("rent===", rent)
}
