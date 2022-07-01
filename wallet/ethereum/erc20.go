package eth

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"
	"savour-core/wallet/ethereum/factory"
)

func (client *EthClient) Erc20BalanceOf(tokenAddress, account string, blockNumber *big.Int) (*big.Int, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}
	tokenContractWrapper, err := factory.NewTokenContractWrapper(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return nil, err
	}
	return tokenContractWrapper.BalanceOfByBlockNumber(nil, common.HexToAddress(account), blockNumber)
}

func (client *EthClient) Erc20Decimals(tokenAddress string) (uint8, error) {
	var decimals uint8
	if client == nil {
		return decimals, errors.New("nil client")
	}
	tokenInstance, err := factory.NewToken(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return decimals, err
	}
	return tokenInstance.Decimals(nil)
}

func (client *EthClient) Erc20RawTransfer(tokenAddress string, nonce uint64, to common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int) (*types.Transaction, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}
	tokenContractWrapper, err := factory.NewTokenContractWrapper(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return nil, err
	}
	return tokenContractWrapper.RawTransfer(
		&bind.TransactOpts{
			Nonce:    big.NewInt(int64(nonce)),
			GasPrice: gasPrice,
			GasLimit: gasLimit,
		},
		to, value)
}

func (client *EthClient) UnpackTransfer(tokenContract common.Address, data []byte) (from, to common.Address, amount *big.Int, err error) {
	if client == nil {
		err = errors.New("nil client")
		return
	}
	tokenContractWrapper, err := factory.NewTokenContractWrapper(tokenContract, client)
	if err != nil {
		fmt.Println(err)
		return
	}
	return tokenContractWrapper.UnpackTransfer(data)
}
