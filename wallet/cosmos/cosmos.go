package cosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/SavourDao/savour-hd/config"
	"github.com/SavourDao/savour-hd/rpc/common"
	"github.com/SavourDao/savour-hd/rpc/wallet"
	wallet2 "github.com/SavourDao/savour-hd/wallet"
	"github.com/SavourDao/savour-hd/wallet/fallback"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/ptypes"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *Client
}

func NewChainAdaptor(conf *config.Config) (wallet2.WalletAdaptor, error) {
	// todo 多个client，按高度来选取
	client, err := NewClient("127.0.0.1:9090")
	if err != nil {
		return nil, err
	}

	return &WalletAdaptor{
		client: client,
	}, nil
}

func (a *WalletAdaptor) Close() error {
	return a.client.Close()
}

func (a *WalletAdaptor) GetBalance(req *wallet.BalanceRequest) (*wallet.BalanceResponse, error) {
	ret, err := a.client.GetBalance(context.Background(), req.Coin, req.Address)
	if err != nil {
		return nil, err
	}

	return &wallet.BalanceResponse{
		Balance: ret.Amount.String(),
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet.TxAddressRequest) (*wallet.TxAddressResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet.TxHashRequest) (*wallet.TxHashResponse, error) {
	ret, err := a.client.GetTxByHash(context.Background(), req.Hash)
	if err != nil {
		return nil, err
	}

	var resp []*Tx
	if err := json.Unmarshal([]byte(ret.TxResponse.Logs.String()), &resp); err != nil {
		return nil, err
	}

	fromAddr, toAddr := "", ""
	for _, v := range resp[0].Events {
		if v.Type == "transfer" {
			for _, attr := range v.Attributes {
				if attr.Key == "recipient" {
					toAddr = attr.Value
				}
				if attr.Key == "sender" {
					fromAddr = attr.Value
				}
			}
		}
	}

	return &wallet.TxHashResponse{
		Tx: &wallet.TxMessage{
			Hash:            req.Hash,
			Index:           uint32(resp[0].MsgIndex),
			Froms:           []*wallet.Address{{Address: fromAddr}},
			Tos:             []*wallet.Address{{Address: toAddr}},
			Values:          nil,
			Fee:             strconv.FormatInt(ret.TxResponse.GasUsed, 10),
			Status:          wallet.TxStatus_Success,
			Type:            0,
			Height:          strconv.FormatInt(ret.TxResponse.Height, 10),
			ContractAddress: "",
			Datetime:        ret.TxResponse.Timestamp,
		},
	}, nil
}

func (a *WalletAdaptor) GetSupportCoins(req *wallet.SupportCoinsRequest) (*wallet.SupportCoinsResponse, error) {
	supportList := []string{"stake", "atom"}

	checkIf := func(s string) bool {
		for _, v := range supportList {
			if strings.EqualFold(v, s) {
				return true
			}
		}
		return false
	}

	return &wallet.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Support: checkIf(req.Chain),
	}, nil
}

func (a *WalletAdaptor) GetNonce(req *wallet.NonceRequest) (*wallet.NonceResponse, error) {
	ret, err := a.client.GetAccount(context.Background(), req.Address)
	if err != nil {
		return nil, err
	}

	account := new(authv1beta1.BaseAccount)
	if err := ptypes.UnmarshalAny(ret.Account, account); err != nil {
		return nil, err
	}

	return &wallet.NonceResponse{
		Nonce: strconv.FormatUint(account.GetSequence(), 10),
	}, nil
}

func (a *WalletAdaptor) GetGasPrice(req *wallet.GasPriceRequest) (*wallet.GasPriceResponse, error) {
	// todo 怎么预估gasfee
	return nil, nil
}

func (a *WalletAdaptor) SendTx(req *wallet.SendTxRequest) (*wallet.SendTxResponse, error) {
	// a.client.SendTx(context.Background(), req.RawTx)
	return nil, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet.ConvertAddressRequest) (*wallet.ConvertAddressResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet.ValidAddressRequest) (*wallet.ValidAddressResponse, error) {
	_, err := types2.AccAddressFromBech32(req.Address)
	if err != nil {
		return &wallet.ValidAddressResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, nil
	}

	return &wallet.ValidAddressResponse{
		Valid: true,
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet.TxFromDataRequest) (*wallet.AccountTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet.TxFromSignedDataRequest) (*wallet.AccountTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet.CreateAccountSignedTxRequest) (*wallet.CreateSignedTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet.CreateAccountTxRequest) (*wallet.CreateAccountTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet.VerifySignedTxRequest) (*wallet.VerifySignedTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet.UtxoInsFromDataRequest) (*wallet.UtxoInsResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet.TxFromDataRequest) (*wallet.UtxoTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet.TxFromSignedDataRequest) (*wallet.UtxoTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet.CreateUtxoSignedTxRequest) (*wallet.CreateSignedTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet.CreateUtxoTxRequest) (*wallet.CreateUtxoTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet.VerifySignedTxRequest) (*wallet.VerifySignedTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet.AccountRequest) (*wallet.AccountResponse, error) {
	ret, err := a.client.GetAccount(context.Background(), req.Address)
	if err != nil {
		return nil, err
	}

	account := new(authv1beta1.BaseAccount)
	if err := ptypes.UnmarshalAny(ret.Account, account); err != nil {
		return nil, err
	}

	return &wallet.AccountResponse{
		AccountNumber: strconv.FormatUint(account.AccountNumber, 10),
		Sequence:      strconv.FormatUint(account.Sequence, 10),
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet.UtxoRequest) (*wallet.UtxoResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet.MinRentRequest) (*wallet.MinRentResponse, error) {
	// todo
	return nil, nil
}
