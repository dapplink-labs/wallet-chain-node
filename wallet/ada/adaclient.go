package ada

import (
	"context"
	"github.com/coinbase/rosetta-sdk-go/asserter"
	rosetta_dk_go_client "github.com/coinbase/rosetta-sdk-go/client"
	"github.com/coinbase/rosetta-sdk-go/types"
	"github.com/pkg/errors"
	"github.com/savour-labs/wallet-hd-chain/config"
	"log"
	"net/http"
)

type adaClient struct {
	client *rosetta_dk_go_client.APIClient
}

func NewAdaClient(conf *config.Config) ([]*adaClient, error) {
	var clients []*adaClient
	for _, rpc := range conf.Fullnode.Ada.RPCs {
		clientCfg := rosetta_dk_go_client.NewConfiguration(
			rpc.RPCURL,
			agent,
			&http.Client{
				Timeout: defaultTimeout,
			},
		)
		client := rosetta_dk_go_client.NewAPIClient(clientCfg)
		aclient := &adaClient{
			client: client,
		}
		clients = append(clients, aclient)
	}
	if len(clients) == 0 {
		return nil, errors.New("No clients available")
	}
	return clients, nil
}

func (c *adaClient) GetUtxoByAddress(address string) ([]*types.Coin, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	accountCoinsRequest := &types.AccountCoinsRequest{
		AccountIdentifier: &types.AccountIdentifier{Address: address},
		NetworkIdentifier: networkIdentifier,
	}
	accountCoinsResponse, rosettaErr, err := c.client.AccountAPI.AccountCoins(ctx, accountCoinsRequest)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	return accountCoinsResponse.Coins, nil
}

func (c *adaClient) GetAccountBalance(address string) ([]*types.Amount, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	accountBalanceRequest := &types.AccountBalanceRequest{
		AccountIdentifier: &types.AccountIdentifier{Address: address},
		NetworkIdentifier: networkIdentifier,
	}
	accountBalanceResponse, rosettaErr, err := c.client.AccountAPI.AccountBalance(
		ctx,
		accountBalanceRequest,
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	return accountBalanceResponse.Balances, nil
}

func (c *adaClient) GetTransactionsByHash(txHash string) ([]*types.BlockTransaction, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	searchTransactionsRequest := &types.SearchTransactionsRequest{
		TransactionIdentifier: &types.TransactionIdentifier{Hash: txHash},
		NetworkIdentifier:     networkIdentifier}
	transactions, rosettaErr, err := c.client.SearchAPI.SearchTransactions(ctx, searchTransactionsRequest)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	return transactions.Transactions, nil
}

func (c *adaClient) GetTransactionsByAddress(address string, limit, offset, maxBlock *int64) ([]*types.BlockTransaction, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	searchTransactionsRequest := &types.SearchTransactionsRequest{
		AccountIdentifier: &types.AccountIdentifier{Address: address},
		Limit:             limit,
		Offset:            offset,
		MaxBlock:          maxBlock,
		NetworkIdentifier: networkIdentifier}
	transactions, rosettaErr, err := c.client.SearchAPI.SearchTransactions(ctx, searchTransactionsRequest)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	return transactions.Transactions, nil
}

func (c *adaClient) GetLatestBlockHeight() (int64, error) {
	ctx := context.Background()
	blockRequest, err := c.GetBlockRequest(ctx)
	if err != nil {
		panic(err)
	}
	blockResponse, rosettaErr, err := c.client.BlockAPI.Block(ctx, blockRequest)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	//log.Printf("Current Block: %s\n", types.PrettyPrintStruct(blockResponse.Block))

	return blockResponse.Block.BlockIdentifier.Index, nil
}

func (c *adaClient) SendRawTransaction(signedTx string) (string, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	txRequest := &types.ConstructionSubmitRequest{
		NetworkIdentifier: networkIdentifier,
		SignedTransaction: signedTx,
	}
	txResponse, rosettaErr, err := c.client.ConstructionAPI.ConstructionSubmit(
		ctx,
		txRequest,
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	return txResponse.TransactionIdentifier.Hash, nil
}

func (c *adaClient) GetNetworkList(ctx context.Context) *types.NetworkIdentifier {
	networkList, rosettaErr, err := c.client.NetworkAPI.NetworkList(
		ctx,
		&types.MetadataRequest{},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}
	if len(networkList.NetworkIdentifiers) == 0 {
		panic("no available networks")
	}
	primaryNetwork := networkList.NetworkIdentifiers[0]
	//log.Printf("Primary Network: %s\n", types.PrettyPrintStruct(networkList))
	return primaryNetwork
}

func (c *adaClient) GetTxFee(relativeTtl, transactionSize int64) ([]*types.Amount, error) {
	ctx := context.Background()
	networkIdentifier := c.GetNetworkList(ctx)
	req := &types.ConstructionMetadataRequest{
		NetworkIdentifier: networkIdentifier,
		Options: map[string]interface{}{
			"relative_ttl":     relativeTtl,
			"transaction_size": transactionSize,
		},
	}
	metaDataResponse, rosettaErr, err := c.client.ConstructionAPI.ConstructionMetadata(ctx, req)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)

	}
	//log.Printf("Rosetta Error: %+v\n", types.PrettyPrintStruct(metaDataResponse))
	return metaDataResponse.SuggestedFee, nil
}

func (c *adaClient) GetBlockRequest(ctx context.Context) (*types.BlockRequest, error) {

	primaryNetwork := c.GetNetworkList(ctx)

	networkStatus, rosettaErr, err := c.client.NetworkAPI.NetworkStatus(
		ctx,
		&types.NetworkRequest{
			NetworkIdentifier: primaryNetwork,
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)

	}

	//log.Printf("Network Status: %s\n", types.PrettyPrintStruct(networkStatus))

	err = asserter.NetworkStatusResponse(networkStatus)
	if err != nil {
		log.Printf("Assertion Error: %s\n", err.Error())
		panic(err)
	}

	networkOptions, rosettaErr, err := c.client.NetworkAPI.NetworkOptions(
		ctx,
		&types.NetworkRequest{
			NetworkIdentifier: primaryNetwork,
		},
	)
	if rosettaErr != nil {
		log.Printf("Rosetta Error: %+v\n", rosettaErr)
		panic(rosettaErr)
	}
	if err != nil {
		panic(err)
	}

	//log.Printf("Network Options: %s\n", types.PrettyPrintStruct(networkOptions))

	err = asserter.NetworkOptionsResponse(networkOptions)
	if err != nil {
		log.Printf("Assertion Error: %s\n", err.Error())
		panic(err)
	}
	return &types.BlockRequest{
		NetworkIdentifier: primaryNetwork,
		BlockIdentifier: types.ConstructPartialBlockIdentifier(
			networkStatus.CurrentBlockIdentifier,
		),
	}, nil
}
