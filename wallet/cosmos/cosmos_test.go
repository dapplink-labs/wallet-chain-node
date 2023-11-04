package cosmos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	wallet2 "github.com/savour-labs/wallet-hd-chain/rpc/wallet"
)

func TestClient_GetBalance(t *testing.T) {
	a, err := NewChainAdaptor(nil)
	assert.NoError(t, err)

	ret, err := a.GetBalance(&wallet2.BalanceRequest{
		Coin:    "stake",
		Address: "cosmos1qd99t24whd3hfg22r53x8uw9ps3rrctwxqvn4m",
	})
	assert.NoError(t, err)

	fmt.Println(ret)
}

func TestClient_GetTxByHash(t *testing.T) {
	a, err := NewChainAdaptor(nil)
	assert.NoError(t, err)

	ret, err := a.GetTxByHash(&wallet2.TxHashRequest{
		Hash: "D43CCA7E9719C0596D766AB4086AC5B02FC74FE49D7CDF368A5FD0386A207396",
	})
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetNonce(t *testing.T) {
	a, err := NewChainAdaptor(nil)
	assert.NoError(t, err)

	ret, err := a.GetNonce(&wallet2.NonceRequest{
		Address: "cosmos1k5s2xws5rlgde0r6kxzs36jlww9ydesgj6vyds",
	})
	assert.NoError(t, err)
	fmt.Println(ret)
}
