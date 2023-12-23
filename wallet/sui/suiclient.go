package sui

import (
	"context"
	"encoding/json"
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

func (c *suiClient) GetTxListByAddress(address string, cursor string, limit uint32) (models.SuiXQueryTransactionBlocksResponse, error) {
	ctx := context.Background()
	req := models.SuiXQueryTransactionBlocksRequest{
		SuiTransactionBlockResponseQuery: models.SuiTransactionBlockResponseQuery{
			TransactionFilter: models.TransactionFilter{
				"FromAddress": address,
			},
			Options: models.SuiTransactionBlockOptions{
				ShowInput:          true,
				ShowRawInput:       true,
				ShowEffects:        true,
				ShowEvents:         true,
				ShowObjectChanges:  true,
				ShowBalanceChanges: true,
			},
		},
		Cursor:          cursor,
		Limit:           uint64(limit),
		DescendingOrder: false,
	}
	txList, err := c.client.SuiXQueryTransactionBlocks(ctx, req)
	if err != nil {
		log.Printf("get tx list  Error: %+v\n", err)
		panic(err)
	}
	return txList, nil
}

func (c *suiClient) GetTxDetailByDigest(digest string) (models.SuiTransactionBlockResponse, error) {
	ctx := context.Background()
	req := models.SuiGetTransactionBlockRequest{
		Digest: digest,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:          true,
			ShowRawInput:       true,
			ShowEffects:        true,
			ShowEvents:         true,
			ShowBalanceChanges: true,
			ShowObjectChanges:  true,
		},
	}
	txDetail, err := c.client.SuiGetTransactionBlock(ctx, req)
	if err != nil {
		log.Printf("get tx detail  Error: %+v\n", err)
		panic(err)
	}
	return txDetail, nil
}

func (c *suiClient) GetGasPrice() (uint64, error) {
	ctx := context.Background()
	price, err := c.client.SuiXGetReferenceGasPrice(ctx)
	if err != nil {
		log.Printf("get gas price Error: %+v\n", err)
		panic(err)
	}
	return price, nil
}

func (c *suiClient) GetCoins(address, coinType string, cursor string, limit uint64) (models.PaginatedCoinsResponse, error) {
	ctx := context.Background()
	if coinType == "" {
		coinType = SuiCoinType
	}
	req := models.SuiXGetCoinsRequest{
		Owner:    address,
		Limit:    limit,
		CoinType: coinType,
	}
	if cursor != "" {
		req.Cursor = cursor
	}

	coins, err := c.client.SuiXGetCoins(ctx, req)
	if err != nil {
		log.Printf("get coins Error: %+v\n", err)
		panic(err)
	}
	return coins, nil
}

func (c *suiClient) GetAllCoins(address string, cursor string,
	limit uint64) (models.PaginatedCoinsResponse, error) {
	ctx := context.Background()
	coinsRequest := models.SuiXGetAllCoinsRequest{
		Owner: address,
		Limit: limit,
	}
	if cursor != "" {
		coinsRequest.Cursor = cursor
	}

	allCoins, err := c.client.SuiXGetAllCoins(ctx, coinsRequest)
	if err != nil {
		log.Printf("get all coins Error: %+v\n", err)
		panic(err)
	}
	return allCoins, nil
}

func (c *suiClient) SendTx(txStr string) (*models.TxnMetaData, error) {
	ctx := context.Background()
	var req models.PublishRequest
	jsonErr := json.Unmarshal([]byte(txStr), &req)
	if jsonErr != nil {
		return nil, jsonErr
	}
	publish, err := c.client.Publish(ctx, req)
	if err != nil {
		log.Printf("publish tx  Error: %+v\n", err)
		panic(err)
	}
	return &publish, nil
}

func (c *suiClient) GetLatestBlockHeight() (int64, error) {
	//TODO implement me
	panic("implement me")
}
