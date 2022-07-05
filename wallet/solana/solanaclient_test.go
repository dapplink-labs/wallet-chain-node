package solana

import (
	"fmt"
	"github.com/SavourDao/savour-core/config"
	"testing"
)

func mockConfig() []*solanaClient {
	client, _ := newSolanaClients(&config.Config{Fullnode: config.Fullnode{
		Sol: config.SolanaNode{
			PublicUrl: "https://public-api.solscan.io",
		},
	},
	})
	return client
}

func TestSolanaClient_GetBalance(t *testing.T) {
	client, _ := newSolanaClients(nil)
	balance := client[0].GetBalance("57vSaRTqN9iXaemgh4AoDsZ63mcaoshfMK8NP3Z5QNbs")
	fmt.Printf("%s", balance)
	if balance != "1.99359136" {
		t.Fatalf("want 1 got %s", balance)
	}
}

func TestSolanaClient_GetTxByHash(t *testing.T) {
	client := mockConfig()
	client[0].GetTxByHash("2DdXMh6QcxnoHnBnAemkzn1QpAYgd1eN1tRxbHQYBmeSYEbjP9yRPA3PNkR7D2HRCcV84oQsYNM2K5KsF9wDBWgE")
}

//func TestSolanaClient_GetFee(t *testing.T) {
//	client, _ := newSolanaClients(nil)
//	balance := client[0].GetFee()
//	fmt.Printf("%s", balance)
//}

func TestSolanaClient_GetTransferHistory(t *testing.T) {
	client, _ := newSolanaClients(nil)
	balance := client[0].GetTxByAddress("57vSaRTqN9iXaemgh4AoDsZ63mcaoshfMK8NP3Z5QNbs")
	fmt.Println(balance)
}

func TestSolanaClient_SendTx(t *testing.T) {
	client := mockConfig()
	client[0].SendTx()
}

func TestSolanaClient_GetAccount(t *testing.T) {
	client := mockConfig()
	address, pri, _ := client[0].GetAccount()
	fmt.Println(address)
	fmt.Println(pri)
}

func TestSolanaClient_RequestAirdrop(t *testing.T) {
	client := mockConfig()
	//client[0].RequestAirdrop("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	balance := client[0].GetBalance("9rZPARQ11UsUcyPDhZ6b98ii4HWYV8wNwfxCBexG8YVX")
	fmt.Println(balance)
}
