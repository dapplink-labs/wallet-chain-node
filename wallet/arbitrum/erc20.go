package arbitrum

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"

	"github.com/savour-labs/wallet-chain-node/wallet/ethereum/factory"
)

func (client *arbiClient) Erc20BalanceOf(tokenAddress, account string, blockNumber *big.Int) (*big.Int, error) {
	if client == nil {
		return nil, errors.New("nil client")
	}
	tokenContractWrapper, err := factory.NewTokenContractWrapper(common.HexToAddress(tokenAddress), client)
	if err != nil {
		return nil, err
	}
	return tokenContractWrapper.BalanceOfByBlockNumber(nil, common.HexToAddress(account), blockNumber)
}

func (client *arbiClient) Erc20Decimals(tokenAddress string) (uint8, error) {
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

func (client *arbiClient) Erc20RawTransfer(tokenAddress string, nonce uint64, to common.Address, value *big.Int, gasLimit uint64, gasPrice *big.Int) (*types.Transaction, error) {
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

func (client *arbiClient) UnpackTransfer(tokenContract common.Address, data []byte) (to common.Address, amount *big.Int, err error) {
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

func (a *WalletAdaptor) validateAndQueryERC20RawTransfer(contract common.Address, tx *types.Transaction) (
	to common.Address, amount *big.Int, err error) {
	// erc20 transfer transaction
	if tx.Value().Cmp(big.NewInt(0)) != 0 {
		log.Error("Invalid ERC20 transfer transaction with non-zero value")
		err = errors.New("Invalid ERC20 transfer transaction with non-zero value")
		return
	}
	if tx.To() == nil || *tx.To() != contract {
		log.Error("Invalid ERC20 transfer with wrong contract address")
		err = errors.New("Invalid ERC20 transfer transaction with wrong contract address")
		return
	}
	if len(tx.Data()) <= 0 {
		log.Error("Invalid ERC20 transfer with empty data")
		err = errors.New("Invalid ERC20 transfer transaction with empty data")
		return
	}
	if len(tx.Data()) <= 0 {
		log.Error("Invalid ERC20 transfer with empty data")
		err = errors.New("Invalid ERC20 transfer transaction with empty data")
		return
	}
	to, amount, err = a.getClient().UnpackTransfer(contract, tx.Data())
	if err != nil {
		log.Error("Invalid ERC20 transfer with unpack error", "err", err)
		err = fmt.Errorf("Invalid ERC20 transfer with unpack error: %v", err)
		return
	}
	return
}

func (a *WalletAdaptor) validateAndQueryERC20TransferReceipt(contract common.Address, from, to string,
	value string, receipt *types.Receipt) error {
	if a.getClient() == nil {
		return errors.New("nil client")
	}
	tokenContractWrapper, err := factory.NewTokenContractWrapper(contract, a.getClient())
	if err != nil {
		log.Error("Failed to create token contract wrapper", "err", err)
		return err
	}

	fromList, toList, valueList, err := tokenContractWrapper.ParseTransferLogs(receipt)
	if err != nil {
		log.Error("Failed to parse transfer logs", "err", err)
		return err
	}
	for i, logFrom := range fromList {
		logTo := toList[i]
		logValue := valueList[i]
		if logFrom.Hex() == from && logTo.Hex() == to &&
			decimal.NewFromBigInt(logValue, 0).String() == value {
			return nil
		}
		log.Error("receipt log does not match transfer",
			"logFrom", logFrom, "logTo", logTo, "logValue", logValue,
			"from", from, "to", to, "value", value)
	}

	log.Error("receipt log does not match transfer",
		"fromList", fromList, "toList", toList, "valueList", valueList)
	return errors.New("receipt log does not match transfer")
}
