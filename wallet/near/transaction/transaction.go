package transaction

import (
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/SavourDao/savour-hd/wallet/near/keys"
	"github.com/near/borsh-go"
)

// Signature asdf.
type Signature struct {
	KeyType uint8
	Data    [64]byte
}

// SignedTransaction asdf.
type SignedTransaction struct {
	Transaction Transaction
	Signature   Signature
}

// Transaction asdf.
type Transaction struct {
	SignerID   string
	PublicKey  PublicKey
	Nonce      uint64
	ReceiverID string
	BlockHash  [32]byte
	Actions    []Action
}

// PublicKey asdf.
type PublicKey struct {
	KeyType uint8
	Data    [32]byte
}

// AccessKey asdf.
type AccessKey struct {
	Nonce      uint64
	Permission AccessKeyPermission
}

// AccessKeyPermission asdf.
type AccessKeyPermission struct {
	Enum         borsh.Enum `borsh_enum:"true"`
	FunctionCall FunctionCallPermission
	FullAccess   FullAccessPermission
}

// FunctionCallPermission asdf.
type FunctionCallPermission struct {
	Allowance   *big.Int
	ReceiverID  string
	MethodNames []string
}

// FullAccessPermission asdf.
type FullAccessPermission struct{}

// Action asdf.
type Action struct {
	Enum           borsh.Enum `borsh_enum:"true"`
	CreateAccount  CreateAccount
	DeployContract DeployContract
	FunctionCall   FunctionCall
	Transfer       Transfer
	Stake          Stake
	AddKey         AddKey
	DeleteKey      DeleteKey
	DeleteAccount  DeleteAccount
}

// CreateAccount asdf.
type CreateAccount struct{}

// DeployContract asdf.
type DeployContract struct {
	Code []byte
}

// FunctionCall asdf.
type FunctionCall struct {
	MethodName string
	Args       []byte
	Gas        uint64
	Deposit    big.Int
}

// Transfer asdf.
type Transfer struct {
	Deposit big.Int
}

// Stake sadf.
type Stake struct {
	Stake     big.Int
	PublicKey PublicKey
}

// AddKey asdf.
type AddKey struct {
	PublicKey PublicKey
	AccessKey AccessKey
}

type DeleteKey struct {
	PublicKey PublicKey
}

// DeleteAccount asdf.
type DeleteAccount struct {
	BeneficiaryID string
}

// TransferAction is a helper to create a Transfer action.
func TransferAction(deposit big.Int) Action {
	return Action{Enum: 3, Transfer: Transfer{Deposit: deposit}}
}

// SignTransaction serializes and signs a Transaction using the provided signer.
func SignTransaction(
	transaction Transaction,
	signer keys.KeyPair,
	accountID string,
	networkID string,
) ([]byte, *SignedTransaction, error) {
	message, err := borsh.Serialize(transaction)
	if err != nil {
		return nil, nil, fmt.Errorf("serializing transaction: %v", err)
	}
	hash := sha256.Sum256(message)
	sig, err := signer.Sign(hash[:])
	if err != nil {
		return nil, nil, fmt.Errorf("signing hash: %v", err)
	}
	var data [64]byte
	copy(data[:], sig)
	st := &SignedTransaction{
		Transaction: transaction,
		Signature: Signature{
			KeyType: transaction.PublicKey.KeyType,
			Data:    data,
		},
	}
	return hash[:], st, nil
}
