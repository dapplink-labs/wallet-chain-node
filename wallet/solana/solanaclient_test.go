package solana

import (
	"fmt"
	"github.com/SavourDao/savour-core/config"
	"testing"
)

func newTestClient() *solanaClient {
	client, _ := newSolanaClients(&config.Config{Fullnode: config.Fullnode{
		Sol: config.SolanaNode{
			PublicUrl: "https://public-api.solscan.io",
			NetWork:   "mainnet",
		},
	},
	})
	return client[0]
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetBalance("8Lh2DVW5Lw3HgmZC55Fquno4K5auzSS7EveuLvEtCEXq")
	fmt.Printf("%s", balance)
}

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := newTestClient()
	client.GetTxByHash("4VtmhbGDANcpFkSsES4sAhmZsdSHQ6MKaijFEuVYbMk35RZziMNs1QLJpZ38HsjVswRJoYBiHNrmUgZigd9exkCJ")
}

func TestSolanaClient_GetTransferHistory(t *testing.T) {
	client := newTestClient()
	balance, _ := client.GetTxByAddress("BqSGA2WdiQXA2cC1EdGDnVD615A4nYEAq49K3fz2hNBo", 1, 10)
	fmt.Println(balance)
}

func TestSolanaClient_GetAccount(t *testing.T) {
	client := newTestClient()
	address, pri, _ := client.GetAccount()
	fmt.Println(address)
	fmt.Println(pri)
}

func TestSolanaClient_RequestAirdrop(t *testing.T) {
	client := newTestClient()
	//client[0].RequestAirdrop("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	balance, _ := client.GetBalance("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	fmt.Println(balance)
}

func TestSolanaClient_SendTx(t *testing.T) {
	client := newTestClient()
	client.SendTx("4voSPg3tYuWbKzimpQK9EbXHmuyy5fUrtXvpLDMLkmY6TRncaTHAKGD8jUg3maB5Jbrd9CkQg4qjJMyN6sQvnEF2")
}

func TestSolanaClient_GetNonce(t *testing.T) {
	client := newTestClient()
	client.GetNonce()
}

func TestSolanaClient_GetMinRent(t *testing.T) {
	client := newTestClient()
	client.GetMinRent()
}
