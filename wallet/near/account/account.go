package account

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/SavourDao/savour-hd/wallet/near/transaction"
	"github.com/SavourDao/savour-hd/wallet/near/types"
	"github.com/mr-tron/base58/base58"
	"github.com/near/borsh-go"
)

// Account provides functions for a single account.
type Account struct {
	config    *types.Config
	accountID string
}

// NewAccount creates a new account.
func NewAccount(config *types.Config, accountID string) *Account {
	return &Account{
		config:    config,
		accountID: accountID,
	}
}

// SignTransaction creates and signs a transaction from the supplied actions.
func (a *Account) SignTransaction(
	ctx context.Context,
	receiverID string,
	accessKey types.GetAccessKeyResponse,
	bHash string,
	actions ...transaction.Action,
) ([]byte, *transaction.SignedTransaction, error) {
	if a.config.Signer == nil {
		return nil, nil, fmt.Errorf("no signer configured")
	}
	blockHash, err := base58.Decode(bHash)
	if err != nil {
		return nil, nil, fmt.Errorf("decoding hash: %v", err)
	}
	var blockHashArr [32]byte
	copy(blockHashArr[:], blockHash)
	nonce := accessKey.Result.Nonce + 1

	pk := a.config.Signer.GetPublicKey()
	var dataArr [32]byte
	copy(dataArr[:], pk.Data)

	t := transaction.Transaction{
		SignerID: a.accountID,
		PublicKey: transaction.PublicKey{
			KeyType: uint8(pk.Type),
			Data:    dataArr,
		},
		Nonce:      uint64(nonce),
		ReceiverID: receiverID,
		BlockHash:  blockHashArr,
		Actions:    actions,
	}
	hash, signedTransaction, err := transaction.SignTransaction(t, a.config.Signer, a.accountID, a.config.NetworkID)
	if err != nil {
		return nil, nil, fmt.Errorf("signing transaction: %v", err)
	}
	return hash, signedTransaction, nil
}

// SignAndSendTransaction creates, signs and sends a tranaction for the supplied actions.
func (a *Account) SignAndSendTransaction(
	ctx context.Context,
	receiverID string,
	accessKey types.GetAccessKeyResponse,
	hash string,
	actions ...transaction.Action,
) (*SendTxResult, error) {
	txHash, signedTransaction, err := a.SignTransaction(ctx, receiverID, accessKey, hash, actions...)
	fmt.Println(txHash)
	if err != nil {
		return nil, err
	}
	bytes, err := borsh.Serialize(*signedTransaction)
	if err != nil {
		return nil, err
	}
	var res SendTxResult
	parms := base64.StdEncoding.EncodeToString(bytes)
	types.DoRpcRequest("broadcast_tx_commit", [1]string{parms}, &res)

	return &res, nil
}
