package tezos

import (
	"fmt"
	"testing"
)

const (
	testScanBaseUrl = "https://api.tzstats.com"
	testScanHash    = "opSrt7oYHDTZcfGnhNt3BzGrrCQf364VuYmKo5ZQVQRfTnczjnf"
)

func TestClient_GetTransactionListByAddress(t *testing.T) {
	result, err := NewTezosScanClient(testScanBaseUrl).GetTransactionByHash(testScanHash)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
