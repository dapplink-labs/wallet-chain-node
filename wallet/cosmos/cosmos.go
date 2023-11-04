package cosmos

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	"github.com/savour-labs/wallet-chain-node/rpc/wallet"
	wallet2 "github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *Client
}

func NewChainAdaptor(conf *config.Config) (wallet2.WalletAdaptor, error) {
	// todo 多个client，按高度来选取
	client, err := NewClient(conf.Fullnode.Cosmos.RPCs[0].RPCURL) // "127.0.0.1:9090"
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
		return &wallet.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}, err
	}

	return &wallet.BalanceResponse{
		Balance: ret.Amount.String(),
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet.TxAddressRequest) (*wallet.TxAddressResponse, error) {
	// todo 接第三方平台接口？支持的第三方机构
	return nil, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet.TxHashRequest) (*wallet.TxHashResponse, error) {
	ret, err := a.client.GetTxByHash(context.Background(), req.Hash)
	if err != nil {
		return &wallet.TxHashResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get tx by hash fail",
		}, err
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

	index := 0
	if resp[0] != nil {
		index = resp[0].MsgIndex
	}

	return &wallet.TxHashResponse{
		Tx: &wallet.TxMessage{
			Hash:            req.Hash,
			Index:           uint32(index),
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
		return &wallet.NonceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get nonce fail",
		}, err
	}

	account := new(authv1beta1.BaseAccount)
	if err := ptypes.UnmarshalAny(ret.Account, account); err != nil {
		return &wallet.NonceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get nonce fail",
		}, err
	}

	return &wallet.NonceResponse{
		Nonce: strconv.FormatUint(account.GetSequence(), 10),
	}, nil
}

func (a *WalletAdaptor) GetGasPrice(req *wallet.GasPriceRequest) (*wallet.GasPriceResponse, error) {
	return &wallet.GasPriceResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) SendTx(req *wallet.SendTxRequest) (*wallet.SendTxResponse, error) {
	// todo test mod
	ret, err := a.client.BroadcastTx(context.Background(), []byte(req.RawTx))
	if err != nil {
		return &wallet.SendTxResponse{
			Code:   common.ReturnCode_ERROR,
			Msg:    "BroadcastTx fail",
			TxHash: ret.TxResponse.TxHash,
		}, err
	}
	return nil, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet.ConvertAddressRequest) (*wallet.ConvertAddressResponse, error) {
	addr := a.client.GetAddressFromPubKey(req.PublicKey)

	return &wallet.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Address: addr,
	}, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet.ValidAddressRequest) (*wallet.ValidAddressResponse, error) {
	_, err := types2.AccAddressFromBech32(req.Address)
	if err != nil {
		return &wallet.ValidAddressResponse{
			Code:  common.ReturnCode_ERROR,
			Msg:   err.Error(),
			Valid: false,
		}, err
	}

	return &wallet.ValidAddressResponse{
		Valid: true,
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet.TxFromDataRequest) (*wallet.AccountTxResponse, error) {
	return &wallet.AccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet.TxFromSignedDataRequest) (*wallet.AccountTxResponse, error) {
	return &wallet.AccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet.CreateAccountSignedTxRequest) (*wallet.CreateSignedTxResponse, error) {
	return &wallet.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet.CreateAccountTxRequest) (*wallet.CreateAccountTxResponse, error) {
	return &wallet.CreateAccountTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet.VerifySignedTxRequest) (*wallet.VerifySignedTxResponse, error) {
	return &wallet.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet.UtxoInsFromDataRequest) (*wallet.UtxoInsResponse, error) {
	return &wallet.UtxoInsResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet.TxFromDataRequest) (*wallet.UtxoTxResponse, error) {
	return &wallet.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet.TxFromSignedDataRequest) (*wallet.UtxoTxResponse, error) {
	return &wallet.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet.CreateUtxoSignedTxRequest) (*wallet.CreateSignedTxResponse, error) {
	return &wallet.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet.CreateUtxoTxRequest) (*wallet.CreateUtxoTxResponse, error) {
	return &wallet.CreateUtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet.VerifySignedTxRequest) (*wallet.VerifySignedTxResponse, error) {
	return &wallet.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
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
	return &wallet.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetMinRent(req *wallet.MinRentRequest) (*wallet.MinRentResponse, error) {
	return &wallet.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}
