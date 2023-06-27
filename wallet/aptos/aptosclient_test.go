package aptos

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetSeq(t *testing.T) {
	c, err := NewClient("https://fullnode.mainnet.aptoslabs.com")
	assert.NoError(t, err)

	ret, err := c.GetSeq("0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12")
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetGasPrice(t *testing.T) {
	c, err := NewClient("https://fullnode.mainnet.aptoslabs.com")
	assert.NoError(t, err)

	ret, err := c.GetGasPrice()
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetTxByAddr(t *testing.T) {
	c, err := NewClient("https://fullnode.mainnet.aptoslabs.com")
	assert.NoError(t, err)

	ret, err := c.GetTxByAddr("0x190d44266241744264b964a37b8f09863167a12d3e70cda39376cfb4e3561e12")
	assert.NoError(t, err)
	fmt.Println(ret[0].Hash)
}

func TestClient_GetTxByTxHash(t *testing.T) {
	c, err := NewClient("https://fullnode.mainnet.aptoslabs.com")
	assert.NoError(t, err)

	ret, err := c.GetTxByTxHash("0xf0274052747da0610e16c1096ef8b08c6056845ac1f1b281b7c465542a8ef37a")
	assert.NoError(t, err)
	fmt.Println(ret)
}

func TestClient_GetTxByVersion(t *testing.T) {
	c, err := NewClient("https://fullnode.mainnet.aptoslabs.com")
	assert.NoError(t, err)

	ret, err := c.GetTxByVersion("123")
	assert.NoError(t, err)
	fmt.Println(ret)
}
