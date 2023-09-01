package filecoin

import (
	"context"
	"strconv"
	"strings"

	"github.com/savour-labs/wallet-hd-chain/config"
	"github.com/savour-labs/wallet-hd-chain/rpc/common"
	"github.com/savour-labs/wallet-hd-chain/rpc/wallet"
	wallet2 "github.com/savour-labs/wallet-hd-chain/wallet"
	"github.com/savour-labs/wallet-hd-chain/wallet/fallback"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *Client
}

func NewChainAdaptor(conf *config.Config) (wallet2.WalletAdaptor, error) {
	// // todo 多个client，按高度来选取
	client, err := NewClient(conf.Fullnode.FileCoin.RPCs[0].RPCURL, conf.Fullnode.FileCoin.ApiToken)
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
	ret, err := a.client.GetBalance(context.Background(), req.Address)
	if err != nil {
		return &wallet.BalanceResponse{
			Code:    common.ReturnCode_ERROR,
			Msg:     "get balance fail",
			Balance: "0",
		}, err
	}

	return &wallet.BalanceResponse{
		Balance: ret.String(),
	}, nil
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet.TxAddressRequest) (*wallet.TxAddressResponse, error) {
	// todo 接第三方平台接口？支持的第三方机构
	return nil, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet.TxHashRequest) (*wallet.TxHashResponse, error) {
	// todo 接第三方平台接口？支持的第三方机构
	return nil, nil
}

func (a *WalletAdaptor) GetSupportCoins(req *wallet.SupportCoinsRequest) (*wallet.SupportCoinsResponse, error) {
	supportList := []string{"fil"}

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
	ret, err := a.client.GetNonce(context.Background(), req.Address)
	if err != nil {
		return &wallet.NonceResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "get nonce fail",
		}, err
	}

	return &wallet.NonceResponse{
		Nonce: strconv.FormatUint(ret, 10),
	}, nil
}

func (a *WalletAdaptor) GetGasPrice(req *wallet.GasPriceRequest) (*wallet.GasPriceResponse, error) {
	return &wallet.GasPriceResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) SendTx(req *wallet.SendTxRequest) (*wallet.SendTxResponse, error) {
	// todo
	return nil, nil
}

func (a *WalletAdaptor) ConvertAddress(req *wallet.ConvertAddressRequest) (*wallet.ConvertAddressResponse, error) {
	addr, err := a.client.GetAddressFromByte(req.PublicKey)
	if err != nil {
		return &wallet.ConvertAddressResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  "fail",
		}, nil
	}

	return &wallet.ConvertAddressResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Address: addr,
	}, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet.ValidAddressRequest) (*wallet.ValidAddressResponse, error) {
	return &wallet.ValidAddressResponse{
		Valid: a.client.VerifyAddress(req.Address),
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
	ret, err := a.client.Send(context.Background(), req.From, req.To, req.Amount)
	if err != nil {
		return &wallet.CreateAccountTxResponse{
			Code: common.ReturnCode_ERROR,
			Msg:  err.Error(),
		}, nil
	}

	return &wallet.CreateAccountTxResponse{
		Code:     common.ReturnCode_SUCCESS,
		TxData:   []byte(ret),
		SignHash: []byte(ret),
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
	return &wallet.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet.UtxoRequest) (*wallet.UtxoResponse, error) {
	return &wallet.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet.MinRentRequest) (*wallet.MinRentResponse, error) {
	return &wallet.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet.ABIBinToJSONRequest) (*wallet.ABIBinToJSONResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet.ABIJSONToBinRequest) (*wallet.ABIJSONToBinResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}
