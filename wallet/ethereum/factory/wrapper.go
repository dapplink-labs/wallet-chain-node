package factory

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// TokenContractWrapper for external wrapping
type TokenContractWrapper struct {
	abi      abi.ABI
	address  common.Address
	backend  bind.ContractBackend
	contract *bind.BoundContract
}

// NewTokenContractWrapper creates TokenContract wrapper for reference
func NewTokenContractWrapper(address common.Address, backend bind.ContractBackend) (*TokenContractWrapper, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	contract, err := bindToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenContractWrapper{
		abi:      parsed,
		address:  address,
		backend:  backend,
		contract: contract,
	}, nil
}

// BalanceOfByBlockNumber with given block number
func (_Token *TokenContractWrapper) BalanceOfByBlockNumber(opts *bind.CallOpts,
	account common.Address, blockNumber *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.call(opts, blockNumber, out, "balanceOf", account)
	return *ret0, err
}

// RawTransfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns()
func (_Token *TokenContractWrapper) RawTransfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error) {
	return _Token.RawTransact(opts, "transfer", _to, _value)
}

// ParseTransferLogs parses log binding the contract event Transfer.
func (_Token *TokenContractWrapper) ParseTransferLogs(receipt *types.Receipt) (
	fromList []common.Address, toList []common.Address, valueList []*big.Int, err error) {
	for _, l := range receipt.Logs {
		event := new(TokenTransfer)
		if len(l.Data) > 0 {
			err = _Token.contract.UnpackLog(event, "Transfer", *l)
			if err == nil {
				fromList = append(fromList, event.From)
				toList = append(toList, event.To)
				valueList = append(valueList, event.Value)
			}
		}
	}
	if len(fromList) <= 0 && err != nil {
		return
	}
	err = nil
	return
}

type transferParams struct {
	To    common.Address
	Value *big.Int
}

type transferFromParams struct {
	From  common.Address
	To    common.Address
	Value *big.Int
}

// UnpackTransfer unpack parameters binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(_to address, _value uint256) returns()
/*func (_Token *TokenContractWrapper) UnpackTransfer(data []byte) (to common.Address, value *big.Int, err error) {
	params := transferParams{}
	method, err := _Token.abi.MethodById(data)
	if err != nil {
		return
	}

	if method.Name != "transfer" {
		err = errors.New(fmt.Sprintf("invalid method, not for transfer %v", method))
		return
	}
	err = _Token.unpackInput(&params, "transfer", data[4:])
	if err != nil {
		return
	}
	to, value = params.To, params.Value
	return
} */

func (_Token *TokenContractWrapper) UnpackTransfer(data []byte) (from, to common.Address, value *big.Int, err error) {
	method, err := _Token.abi.MethodById(data)
	if err != nil {
		return
	}
	switch method.Name {
	case "transfer":
		params := transferParams{}
		err = _Token.unpackInput(&params, "transfer", data[4:])
		if err != nil {
			return
		}
		from = common.Address{}
		to, value = params.To, params.Value
	case "transferFrom":
		params := transferFromParams{}
		err = _Token.unpackInput(&params, "transferFrom", data[4:])
		if err != nil {
			return
		}
		from, to, value = params.From, params.To, params.Value
	default:
		err = errors.New("invalid method, not for transfer")
		return
	}
	return
}

// call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
// from go-ethereum/accounts/abi/bind/base.go
func (_Token *TokenContractWrapper) call(opts *bind.CallOpts, blockNumber *big.Int,
	result interface{}, method string, params ...interface{}) error {
	// Don't crash on a lazy user
	if opts == nil {
		opts = new(bind.CallOpts)
	}
	// Pack the input, call and unpack the results
	input, err := _Token.abi.Pack(method, params...)
	if err != nil {
		return err
	}
	var (
		msg    = ethereum.CallMsg{From: opts.From, To: &_Token.address, Data: input}
		ctx    = ensureContext(opts.Context)
		code   []byte
		output []byte
	)
	if opts.Pending {
		pb, ok := _Token.backend.(bind.PendingContractCaller)
		if !ok {
			return bind.ErrNoPendingState
		}
		output, err = pb.PendingCallContract(ctx, msg)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = pb.PendingCodeAt(ctx, _Token.address); err != nil {
				return err
			} else if len(code) == 0 {
				return bind.ErrNoCode
			}
		}
	} else {
		output, err = _Token.backend.CallContract(ctx, msg, blockNumber)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = _Token.backend.CodeAt(ctx, _Token.address, nil); err != nil {
				return err
			} else if len(code) == 0 {
				return bind.ErrNoCode
			}
		}
	}
	if err != nil {
		return err
	}
	return _Token.abi.Unpack(result, method, output)
}

// RawTransact invokes the (paid) contract method with params as input values.
// from go-ethereum/accounts/abi/bind/base.go
// modification: do not sign or send the transaction
func (_Token *TokenContractWrapper) RawTransact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	// Otherwise pack up the parameters and invoke the contract
	input, err := _Token.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}
	return _Token.rawTransact(opts, &_Token.address, input)
}

// rawTransact executes an actual transaction invocation, first deriving any missing
// authorization fields, and then scheduling the transaction for execution.
// from go-ethereum/accounts/abi/bind/base.go
// modification: do not sign or send the transaction
func (_Token *TokenContractWrapper) rawTransact(opts *bind.TransactOpts, contract *common.Address, input []byte) (*types.Transaction, error) {
	// Ensure a valid value field and resolve the account nonce
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}
	var nonce uint64
	if opts.Nonce == nil {
		return nil, errors.New("failed to retrieve account nonce")
	}

	nonce = opts.Nonce.Uint64()
	// Figure out the gas allowance and gas price values
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		return nil, errors.New("failed to suggest gas price")
	}
	gasLimit := opts.GasLimit
	if gasLimit == 0 {
		return nil, errors.New("failed to estimate gas needed: %v")
	}
	// Create the transaction, sign it and schedule it for execution
	var rawTx *types.Transaction
	if contract == nil {
		rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	} else {
		rawTx = types.NewTransaction(nonce, _Token.address, value, gasLimit, gasPrice, input)
	}
	return rawTx, nil
}

// unpackInput input in v according to the abi specification
// from Unpack() in go-ethereum/accounts/abi/abi.go
func (_Token *TokenContractWrapper) unpackInput(v interface{}, name string, input []byte) (err error) {
	if len(input) == 0 {
		return fmt.Errorf("abi: unmarshalling empty input")
	}
	// since there can't be naming collisions with contracts and events,
	// we need to decide whether we're calling a method or an event
	if method, ok := _Token.abi.Methods[name]; ok {
		if len(input)%32 != 0 {
			return fmt.Errorf("abi: improperly formatted output: %s - Bytes: [%+v]", string(input), input)
		}
		return method.Inputs.Unpack(v, input)
	}
	return fmt.Errorf("abi: could not locate named method")
}

// ensureContext is a helper method to ensure a context is not nil, even if the
// user specified it as such.
// from go-ethereum/accounts/abi/bind/base.go
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
