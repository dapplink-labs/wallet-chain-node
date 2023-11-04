package sui

import (
	"context"
	"github.com/block-vision/sui-go-sdk/utils"
	"log"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/pkg/errors"

	"github.com/savour-labs/wallet-chain-node/config"
)

type suiClient struct {
	client sui.ISuiAPI
}

func NewSuiClient(conf *config.Config) ([]*suiClient, error) {
	var clients []*suiClient
	for _, rpc := range conf.Fullnode.Sui.RPCs {
		client := sui.NewSuiClient(rpc.RPCURL)
		aclient := &suiClient{
			client: client,
		}
		clients = append(clients, aclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (c *suiClient) GetAccountBalance(owner, coinType string) (models.CoinBalanceResponse, error) {
	ctx := context.Background()
	// if coinType is empty, use default coin type
	if coinType == "" {
		coinType = SuiCoinType
	}
	req := models.SuiXGetBalanceRequest{
		Owner:    owner,
		CoinType: coinType,
	}
	balance, err := c.client.SuiXGetBalance(ctx, req)
	if err != nil {
		log.Printf("get balance Error: %+v\n", err)
		panic(err)
	}
	return balance, nil
}

func (c *suiClient) GetAllAccountBalance(owner string) (models.CoinAllBalanceResponse, error) {
	ctx := context.Background()
	req := models.SuiXGetAllBalanceRequest{
		Owner: owner,
	}
	balance, err := c.client.SuiXGetAllBalance(ctx, req)
	if err != nil {
		log.Printf("get all balance Error: %+v\n", err)
		panic(err)
	}
	return balance, nil
}

func (c *suiClient) GetTxListByAddress(address string) {
	ctx := context.Background()
	blocks, err := c.client.SuiGetTotalTransactionBlocks(ctx)
	if err != nil {
		log.Printf("get all balance Error: %+v\n", err)
		panic(err)
	}
	utils.PrettyPrint(blocks)
}

func (c *suiClient) GetTxDetailByHash(hash string) {

}

func (c *suiClient) GetTxFee() {

}

func (c *suiClient) SendTx(signedTx string) (string, error) {
	panic("not implemented")
}
