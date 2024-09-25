package bitcoin

import (
	"errors"
	"fmt"
	"testing"
)

const (
	testBaseUrl = "https://blockchain.info/"
	testAddress = "bc1ql49ydapnjafl5t2cp9zqpjwe6pdgmxy98859v2"
)

var errBlockChainURLEmpty = errors.New("blockchain URL cannot be empty")

func TestNewBlockChainClient_EmptyURL(t *testing.T) {
	url := ""

	client, err := NewBlockChainClient(url)

	if client != nil {
		t.Fatalf("Expected client to be nil when URL is empty, but got %v", client)
	}

	if err == nil || err.Error() != errBlockChainURLEmpty.Error() {
		t.Fatalf("Expected error: %v, but got: %v", errBlockChainURLEmpty, err)
	}
}

func TestBcClient_GetAccountBalance(t *testing.T) {
	client, err := NewBlockChainClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	balanceA, err := client.GetAccountBalance(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balanceA)
}

func TestBcClient_GetAccountUtxo(t *testing.T) {
	client, err := NewBlockChainClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	utxoList, err := client.GetAccountUtxo(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(utxoList)
}

func TestBcClient_GetTransactionsByAddress(t *testing.T) {
	client, err := NewBlockChainClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	txList, err := client.GetTransactionsByAddress(testAddress, "10", "1")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(txList)
}

func TestBcClient_GetTransactionsByHash(t *testing.T) {
	client, err := NewBlockChainClient(testBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	transaction, err := client.GetTransactionsByHash("967fa0c625c565d50032cadd2aba3a564160f0008dfac663b70584abf55d551b")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(transaction)
}
