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
	ret, err := NewTezosClient(testBaseUrl).GetAccountBalance(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}

func TestClient_GetAccountCounter(t *testing.T) {
	ret, err := NewTezosClient(testBaseUrl).GetAccountCounter(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}

func TestClient_GetManageKey(t *testing.T) {
	ret, err := NewTezosClient(testBaseUrl).GetManagerKey(testAddress)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ret)
}
