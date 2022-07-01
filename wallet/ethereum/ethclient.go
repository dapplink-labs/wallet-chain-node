package eth

import (
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type EthClient struct {
	Client
	chainConfig   *params.ChainConfig
	confirmations uint64
}

type Client interface {
	bind.ContractBackend
	BalanceAt(context.Context, common.Address, *big.Int) (*big.Int, error)
	TransactionByHash(context.Context, common.Hash) (*types.Transaction, bool, error)
	BlockByNumber(context.Context, *big.Int) (*types.Block, error)
	TransactionReceipt(context.Context, common.Hash) (*types.Receipt, error)
	NonceAt(context.Context, common.Address, *big.Int) (uint64, error)
}

func NewEthClients(rpc_url string) (*EthClient, error) {
	chainConfig := params.MainnetChainConfig
	client := &EthClient{
		chainConfig:   chainConfig,
		confirmations: 1,
	}
	var err error
	client.Client, err = ethclient.Dial(rpc_url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (this EthClient) GetEthBalance(address string) (*big.Int, error) {
	balance, err := this.BalanceAt(context.TODO(), common.HexToAddress(address), nil)
	return balance, err
}

func (this EthClient) GetNonce(address string) (uint64, error) {
	nonce, err := this.NonceAt(context.TODO(), common.HexToAddress(address), nil)
	return nonce, err
}

func (this EthClient) GetGasPrice(address string) (*big.Int, error) {
	gas, err := this.SuggestGasPrice(context.TODO())
	return gas, err
}

func (this EthClient) SendTxHex(strHex string) (string, error) {
	txbytes, err := hexutil.Decode(strHex)
	if err != nil {
		return "", err
	}
	var txsigned types.Transaction
	err = rlp.DecodeBytes(txbytes, &txsigned)
	if err != nil {
		return "", err
	}
	err = this.SendTransaction(context.Background(), &txsigned)
	if err != nil {
		return "", err
	}
	return txsigned.Hash().String(), nil
}
