package cosmos

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetTxByEvent(t *testing.T) {
	c, err := NewClient("127.0.0.1:9090")
	assert.NoError(t, err)

	event := []string{"send"}
	ret, err := c.GetTxByEvent(context.Background(), event, 0, 0)
	assert.NoError(t, err)
	fmt.Println("here", ret)
}

func TestClient_GetAccount(t *testing.T) {
	c, err := NewClient("127.0.0.1:9090")
	assert.NoError(t, err)

	ret, err := c.GetAccount(context.Background(), "cosmos1qd99t24whd3hfg22r53x8uw9ps3rrctwxqvn4m")
	assert.NoError(t, err)
	fmt.Println("here", ret)
}
