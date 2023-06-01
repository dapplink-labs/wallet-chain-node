package tezos

import (
	"fmt"
	"testing"
)

const (
	testBaseUrl = "https://mainnet.api.tez.ie"
	testAddress = "tz1aLrL64dyofDJQSP8rBip9GykihykWX548"
)

func TestClient_GetAccountBalance(t *testing.T) {
	ret, err := NewTezosClient(testBaseUrl).getAccountBalance(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}

func TestClient_GetAccountCounter(t *testing.T) {
	ret, err := NewTezosClient(testBaseUrl).getAccountCounter(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}

func TestClient_GetManageKey(t *testing.T) {
	ret, err := NewTezosClient(testBaseUrl).getManagerKey(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}
