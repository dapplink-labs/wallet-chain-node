package tron

import (
	"fmt"
	"testing"
)

const (
	testScanBaseUrl = "https://api.trongrid.io"
	testScanHash    = "opSrt7oYHDTZcfGnhNt3BzGrrCQf364VuYmKo5ZQVQRfTnczjnf"
)

func TestClient_GetBalance(t *testing.T) {
	client, err := NewTronGridClient(testScanBaseUrl)
	if err != nil {
		fmt.Println(err)
	}
	balanceData := client.GetBalance("")
	fmt.Println(balanceData)
}
