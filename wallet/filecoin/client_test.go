package filecoin

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetBalancet(t *testing.T) { // addr: f3r5rnril6w773yfrrix7h4ulx6rg3ky2m6c4aae5euxta23lcktzo5ede3gogwcxouo6iie5fpop2l4hbkoqq
	c, err := NewClient("127.0.0.1:1234", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.5eY7Nokrw0eKn6kah21RH2lawg6rbHN1l9mw2s2YxRA")
	assert.NoError(t, err)

	// newAdd, err := c.GetWallet(context.Background())
	// assert.NoError(t, err)
	// fmt.Println(newAdd.String())
	//
	// addr, err := c.GetAddress(context.Background())
	// assert.NoError(t, err)
	// fmt.Println(addr.String())

	ret, err := c.GetBalance(context.Background(), "f3r5rnril6w773yfrrix7h4ulx6rg3ky2m6c4aae5euxta23lcktzo5ede3gogwcxouo6iie5fpop2l4hbkoqq")
	assert.NoError(t, err)
	fmt.Println(ret)

	// ret1, err := c.GetNonce(context.Background(), "f3r5rnril6w773yfrrix7h4ulx6rg3ky2m6c4aae5euxta23lcktzo5ede3gogwcxouo6iie5fpop2l4hbkoqq")
	// assert.NoError(t, err)
	// fmt.Println(ret1)

	ret2, err := c.Send(context.Background(), "f3rtelnhtvvgnnvyb2zrp5xi3x7632p5g24q2djot6jtr4ufx3ysnfdhfdwzjvfxknt2kmsw4jdyhxqodo5ohq", "f3r5rnril6w773yfrrix7h4ulx6rg3ky2m6c4aae5euxta23lcktzo5ede3gogwcxouo6iie5fpop2l4hbkoqq", "1")
	assert.NoError(t, err)
	fmt.Println(ret2)
}
