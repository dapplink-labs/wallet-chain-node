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
			NetWork:   "devnet",
		},
	},
	})
	return client[0]
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client := newTestClient()
	balance := client.GetBalance("RNfp4xTbBb4C3kcv2KqtAj8mu4YhMHxqm1Skg9uchZ7")
	fmt.Printf("%s", balance)
}

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := newTestClient()
	client.GetTxByHash("4F78cfDYrddKwHxQUrkBurGNMh6xfNee47XNbN91Wn7jNQmV3yhTC9ELBMA91FFJTbGtovYeXgPNmPYSJkFd4C8v")
}

func TestSolanaClient_GetTransferHistory(t *testing.T) {
	client, _ := newSolanaClients(nil)
	balance := client[0].GetTxByAddress("57vSaRTqN9iXaemgh4AoDsZ63mcaoshfMK8NP3Z5QNbs")
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
	balance := client.GetBalance("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	fmt.Println(balance)
}

func TestSolanaClient_SendTx(t *testing.T) {
	client := newTestClient()
	client.SendTx()
}

func TestSolanaClient_GetNonce(t *testing.T) {
	client := newTestClient()
	client.GetNonce("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
}

func TestSolanaClient_GetMinRent(t *testing.T) {
	client := newTestClient()
	client.GetMinRent()
}
