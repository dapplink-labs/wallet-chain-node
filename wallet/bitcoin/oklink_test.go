package bitcoin

import (
	"fmt"
	"testing"
)

func TestOkLinkClient_GetGasFee(t *testing.T) {
	okLinkClient, err := NewOkLinkClient("https://www.oklink.com")
	if err != nil {
		fmt.Println(err)
	}
	balanceA, err := okLinkClient.GetGasFee("btc")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(balanceA)
}
