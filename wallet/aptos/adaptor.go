package aptos

import (
	"strconv"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/savour-labs/wallet-chain-node/config"
	"github.com/savour-labs/wallet-chain-node/rpc/common"
	wallet2 "github.com/savour-labs/wallet-chain-node/rpc/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet"
	"github.com/savour-labs/wallet-chain-node/wallet/fallback"
)

type WalletAdaptor struct {
	fallback.WalletAdaptor
	client *Client
}

func NewChainAdaptor(conf *config.Config) (wallet.WalletAdaptor, error) {
	client, err := NewClient(conf.Fullnode.Oasis.RPCs[0].RPCURL)
	if err != nil {
		return nil, err
	}
	return &WalletAdaptor{
		client: client,
	}, nil
}

func (a *WalletAdaptor) GetBalance(req *wallet2.BalanceRequest) (*wallet2.BalanceResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetTxByAddress(req *wallet2.TxAddressRequest) (*wallet2.TxAddressResponse, error) {
	ret, err := a.client.GetTxByAddr(req.Address)
	if err != nil {
		return nil, err
	}

	res := &wallet2.TxAddressResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
	}

	res.Tx = make([]*wallet2.TxMessage, 0)
	for _, v := range ret {
		index, err := strconv.Atoi(v.SequenceNumber)
		if err != nil {
			return nil, err
		}
		res.Tx = append(res.Tx, &wallet2.TxMessage{
			Hash:  v.Hash,
			Index: uint32(index),
		})
	}

	return res, nil
}

func (a *WalletAdaptor) GetTxByHash(req *wallet2.TxHashRequest) (*wallet2.TxHashResponse, error) {
	ret, err := a.client.GetTxByTxHash(req.Hash)
	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(ret.SequenceNumber)
	if err != nil {
		return nil, err
	}
	res := &wallet2.TxHashResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Tx: &wallet2.TxMessage{
			Hash:  ret.Hash,
			Index: uint32(index),
		},
	}

	return res, nil
}

func (a *WalletAdaptor) GetSupportCoins(req *wallet2.SupportCoinsRequest) (*wallet2.SupportCoinsResponse, error) {
	return &wallet2.SupportCoinsResponse{
		Code:    common.ReturnCode_SUCCESS,
		Msg:     "success",
		Support: true,
	}, nil
}

func (a *WalletAdaptor) GetNonce(req *wallet2.NonceRequest) (*wallet2.NonceResponse, error) {
	ret, err := a.client.GetSeq(req.Address)
	if err != nil {
		return nil, err
	}

	return &wallet2.NonceResponse{
		Code:  common.ReturnCode_SUCCESS,
		Msg:   "success",
		Nonce: strconv.Itoa(ret),
	}, nil
}

func (a *WalletAdaptor) GetGasPrice(req *wallet2.GasPriceRequest) (*wallet2.GasPriceResponse, error) {
	ret, err := a.client.GetGasPrice()
	if err != nil {
		return nil, err
	}

	return &wallet2.GasPriceResponse{
		Code: common.ReturnCode_SUCCESS,
		Msg:  "success",
		Gas:  strconv.Itoa(ret.GasEstimate),
	}, nil
}

func (a *WalletAdaptor) SendTx(req *wallet2.SendTxRequest) (*wallet2.SendTxResponse, error) {
	// todo api需要的信息不足
	// TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ConvertAddress(req *wallet2.ConvertAddressRequest) (*wallet2.ConvertAddressResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) ValidAddress(req *wallet2.ValidAddressRequest) (*wallet2.ValidAddressResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetAccountTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetAccountTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountSignedTx(req *wallet2.CreateAccountSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) CreateAccountTx(req *wallet2.CreateAccountTxRequest) (*wallet2.CreateAccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) VerifyAccountSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return nil, nil
}

// queryTransaction retrieve transaction information from a signed data.
func (a *WalletAdaptor) queryTransaction(isERC20 bool, tx *types.Transaction, receipt *types.Receipt, blockNumber uint64, signer types.Signer) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

// queryRawTransaction retrieve transaction information from a raw(unsigned) data.
func (a *WalletAdaptor) queryRawTransaction(isERC20 bool, rawTx *types.Transaction) (*wallet2.AccountTxResponse, error) {
	return nil, nil
}

func (a *WalletAdaptor) GetUtxoInsFromData(req *wallet2.UtxoInsFromDataRequest) (*wallet2.UtxoInsResponse, error) {
	return &wallet2.UtxoInsResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromData(req *wallet2.TxFromDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetUtxoTxFromSignedData(req *wallet2.TxFromSignedDataRequest) (*wallet2.UtxoTxResponse, error) {
	return &wallet2.UtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoSignedTx(req *wallet2.CreateUtxoSignedTxRequest) (*wallet2.CreateSignedTxResponse, error) {
	return &wallet2.CreateSignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) CreateUtxoTx(req *wallet2.CreateUtxoTxRequest) (*wallet2.CreateUtxoTxResponse, error) {
	return &wallet2.CreateUtxoTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) VerifyUtxoSignedTx(req *wallet2.VerifySignedTxRequest) (*wallet2.VerifySignedTxResponse, error) {
	return &wallet2.VerifySignedTxResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Don't support",
	}, nil
}

func (a *WalletAdaptor) GetAccount(req *wallet2.AccountRequest) (*wallet2.AccountResponse, error) {
	return &wallet2.AccountResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetUtxo(req *wallet2.UtxoRequest) (*wallet2.UtxoResponse, error) {
	return &wallet2.UtxoResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) GetMinRent(req *wallet2.MinRentRequest) (*wallet2.MinRentResponse, error) {
	return &wallet2.MinRentResponse{
		Code: common.ReturnCode_ERROR,
		Msg:  "Do not support this interface",
	}, nil
}

func (a *WalletAdaptor) ABIBinToJSON(req *wallet2.ABIBinToJSONRequest) (*wallet2.ABIBinToJSONResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) ABIJSONToBin(req *wallet2.ABIJSONToBinRequest) (*wallet2.ABIJSONToBinResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (a *WalletAdaptor) GetUnspentOutputs(req *wallet2.UnspentOutputsRequest) (*wallet2.UnspentOutputsResponse, error) {
	//TODO implement me
	panic("implement me")
}
